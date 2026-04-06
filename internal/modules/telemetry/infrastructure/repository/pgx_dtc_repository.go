package telemetryrepository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type PgxDTCRepository struct {
	pool *pgxpool.Pool
}

func NewPgxDTCRepository(pool *pgxpool.Pool) *PgxDTCRepository {
	return &PgxDTCRepository{pool: pool}
}

func (r *PgxDTCRepository) Insert(ctx context.Context, entry *telemetrydomain.DTCEntry) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO dtc_entries (time, device_id, vin, code, status) VALUES ($1, $2, $3, $4, $5)`,
		entry.Time, entry.DeviceID, entry.VIN, entry.Code, entry.Status,
	)
	return err
}

func (r *PgxDTCRepository) ListActive(ctx context.Context, vin string) ([]*telemetrydomain.DTCEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (code) time, device_id, vin, code, status
		FROM dtc_entries
		WHERE vin = $1
		ORDER BY code, time DESC
	`, vin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*telemetrydomain.DTCEntry
	for rows.Next() {
		e := &telemetrydomain.DTCEntry{}
		if err := rows.Scan(&e.Time, &e.DeviceID, &e.VIN, &e.Code, &e.Status); err != nil {
			return nil, err
		}
		if e.Status == "opened" {
			entries = append(entries, e)
		}
	}
	return entries, rows.Err()
}

func (r *PgxDTCRepository) HasActiveDTCs(ctx context.Context, vins []string) (map[string]bool, error) {
	if len(vins) == 0 {
		return map[string]bool{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT vin, status
		FROM (
			SELECT DISTINCT ON (vin, code) vin, status
			FROM dtc_entries
			WHERE vin = ANY($1::text[])
			ORDER BY vin, code, time DESC
		) latest
		WHERE status = 'opened'
	`, vins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]bool, len(vins))
	for rows.Next() {
		var vin, status string
		if err := rows.Scan(&vin, &status); err != nil {
			return nil, err
		}
		result[vin] = true
	}
	return result, rows.Err()
}
