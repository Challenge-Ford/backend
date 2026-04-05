package telemetryrepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type PgxDTCRepository struct {
	pool *pgxpool.Pool
}

func NewPgxDTCRepository(pool *pgxpool.Pool) *PgxDTCRepository {
	return &PgxDTCRepository{pool: pool}
}

func (r *PgxDTCRepository) Save(ctx context.Context, dtc *telemetrydomain.ActiveDTC) error {
	if dtc.IsClosed() {
		_, err := r.pool.Exec(ctx,
			`DELETE FROM telemetry.active_dtcs WHERE device_id = $1 AND code = $2`,
			dtc.DeviceID, dtc.Code,
		)
		return err
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO telemetry.active_dtcs (device_id, vin, code, detected_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (device_id, code) DO UPDATE SET vin = EXCLUDED.vin, detected_at = EXCLUDED.detected_at`,
		dtc.DeviceID, dtc.VIN, dtc.Code, dtc.DetectedAt,
	)
	return err
}

func (r *PgxDTCRepository) ListActive(ctx context.Context, vin string) ([]*telemetrydomain.ActiveDTC, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT device_id, vin, code, detected_at FROM telemetry.active_dtcs WHERE vin = $1 ORDER BY detected_at DESC`,
		vin,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dtcs []*telemetrydomain.ActiveDTC
	for rows.Next() {
		var deviceID uuid.UUID
		d := &telemetrydomain.ActiveDTC{}
		if err := rows.Scan(&deviceID, &d.VIN, &d.Code, &d.DetectedAt); err != nil {
			return nil, err
		}
		d.DeviceID = deviceID
		dtcs = append(dtcs, d)
	}
	return dtcs, rows.Err()
}

func (r *PgxDTCRepository) HasActiveDTCs(ctx context.Context, vins []string) (map[string]bool, error) {
	if len(vins) == 0 {
		return map[string]bool{}, nil
	}
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT vin FROM telemetry.active_dtcs WHERE vin = ANY($1)`,
		vins,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]bool, len(vins))
	for rows.Next() {
		var vin string
		if err := rows.Scan(&vin); err != nil {
			return nil, err
		}
		result[vin] = true
	}
	return result, rows.Err()
}
