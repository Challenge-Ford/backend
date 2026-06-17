package telemetryrepository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type PgxStateObservationRepository struct {
	pool *pgxpool.Pool
}

func NewPgxStateObservationRepository(pool *pgxpool.Pool) *PgxStateObservationRepository {
	return &PgxStateObservationRepository{pool: pool}
}

func (r *PgxStateObservationRepository) Insert(ctx context.Context, e *telemetrydomain.VehicleStateObservation) (bool, error) {
	state, err := json.Marshal(e.State)
	if err != nil {
		return false, err
	}
	observation, err := json.Marshal(e.Observation)
	if err != nil {
		return false, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		INSERT INTO vehicle_state_message_ids (message_id, observed_at, received_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (message_id) DO NOTHING`,
		e.MessageID, e.ObservedAt, e.ReceivedAt,
	)
	if err != nil {
		return false, err
	}
	if tag.RowsAffected() == 0 {
		return false, nil
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO vehicle_state_observations (
			observed_at, message_id, device_id, vehicle_id, received_at,
			state, observation, raw_payload
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7::jsonb, $8::jsonb)
		ON CONFLICT (observed_at, message_id) DO NOTHING`,
		e.ObservedAt, e.MessageID, e.DeviceID, e.VehicleID, e.ReceivedAt,
		state, observation, e.RawPayload,
	); err != nil {
		return false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (r *PgxStateObservationRepository) List(ctx context.Context, vehicleID uuid.UUID, from, to *time.Time, limit int, after *time.Time) ([]*telemetrydomain.VehicleStateObservation, error) {
	args := []any{vehicleID}
	q := `
		SELECT observed_at, message_id, device_id, vehicle_id, received_at, state, observation, raw_payload
		FROM vehicle_state_observations
		WHERE vehicle_id = $1`
	argIdx := 2
	if from != nil {
		q += ` AND observed_at >= $` + strconv.Itoa(argIdx)
		args = append(args, *from)
		argIdx++
	}
	if to != nil {
		q += ` AND observed_at <= $` + strconv.Itoa(argIdx)
		args = append(args, *to)
		argIdx++
	}
	if after != nil {
		q += ` AND observed_at > $` + strconv.Itoa(argIdx)
		args = append(args, *after)
		argIdx++
	}
	q += ` ORDER BY observed_at ASC LIMIT $` + strconv.Itoa(argIdx)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*telemetrydomain.VehicleStateObservation
	for rows.Next() {
		e, err := scanObservation(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *PgxStateObservationRepository) Latest(ctx context.Context, vehicleID uuid.UUID) (*telemetrydomain.VehicleStateObservation, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT observed_at, message_id, device_id, vehicle_id, received_at, state, observation, raw_payload
		FROM vehicle_state_observations
		WHERE vehicle_id = $1
		ORDER BY observed_at DESC
		LIMIT 1`, vehicleID)

	e, err := scanObservation(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *PgxStateObservationRepository) LatestPosition(ctx context.Context, vehicleID uuid.UUID) (*telemetrydomain.VehicleStateObservation, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT observed_at, message_id, device_id, vehicle_id, received_at, state, observation, raw_payload
		FROM vehicle_state_observations
		WHERE vehicle_id = $1
		  AND state ? 'position'
		ORDER BY observed_at DESC
		LIMIT 1`, vehicleID)

	e, err := scanObservation(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *PgxStateObservationRepository) ListActiveDTCs(ctx context.Context, vehicleID uuid.UUID) ([]*telemetrydomain.ActiveDTC, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT observed_at, state
		FROM vehicle_state_observations
		WHERE vehicle_id = $1
		  AND state ? 'diagnostics'
		  AND state->'diagnostics' ? 'open_dtcs'
		ORDER BY observed_at DESC
		LIMIT 1`, vehicleID)

	var observedAt time.Time
	var rawState []byte
	if err := row.Scan(&observedAt, &rawState); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*telemetrydomain.ActiveDTC{}, nil
		}
		return nil, err
	}
	var state telemetrydomain.VehicleState
	if err := json.Unmarshal(rawState, &state); err != nil {
		return nil, err
	}
	if state.Diagnostics == nil || state.Diagnostics.OpenDTCs == nil {
		return []*telemetrydomain.ActiveDTC{}, nil
	}
	out := make([]*telemetrydomain.ActiveDTC, 0, len(state.Diagnostics.OpenDTCs))
	for _, code := range state.Diagnostics.OpenDTCs {
		out = append(out, &telemetrydomain.ActiveDTC{Code: code, Time: observedAt})
	}
	return out, nil
}

func (r *PgxStateObservationRepository) HasActiveDTCs(ctx context.Context, vehicleIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	if len(vehicleIDs) == 0 {
		return map[uuid.UUID]bool{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (vehicle_id)
			vehicle_id,
			jsonb_array_length(state->'diagnostics'->'open_dtcs') > 0 AS has_active_dtcs
		FROM vehicle_state_observations
		WHERE vehicle_id = ANY($1::uuid[])
		  AND state ? 'diagnostics'
		  AND state->'diagnostics' ? 'open_dtcs'
		ORDER BY vehicle_id, observed_at DESC`, vehicleIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]bool, len(vehicleIDs))
	for rows.Next() {
		var vehicleID uuid.UUID
		var hasActive bool
		if err := rows.Scan(&vehicleID, &hasActive); err != nil {
			return nil, err
		}
		if hasActive {
			result[vehicleID] = true
		}
	}
	return result, rows.Err()
}

type observationScanner interface {
	Scan(dest ...any) error
}

func scanObservation(s observationScanner) (*telemetrydomain.VehicleStateObservation, error) {
	e := &telemetrydomain.VehicleStateObservation{}
	var rawState []byte
	var rawObservation []byte
	if err := s.Scan(
		&e.ObservedAt, &e.MessageID, &e.DeviceID, &e.VehicleID, &e.ReceivedAt,
		&rawState, &rawObservation, &e.RawPayload,
	); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rawState, &e.State); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rawObservation, &e.Observation); err != nil {
		return nil, err
	}
	return e, nil
}
