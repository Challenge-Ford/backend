package telemetryrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type PgxDTCatalogRepository struct {
	pool *pgxpool.Pool
}

func NewPgxDTCatalogRepository(pool *pgxpool.Pool) *PgxDTCatalogRepository {
	return &PgxDTCatalogRepository{pool: pool}
}

func (r *PgxDTCatalogRepository) GetWithEstimates(ctx context.Context, code string, modelYearID uuid.UUID) (*telemetrydomain.DTCCatalogWithEstimates, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT c.code, c.description, c.system, c.severity, c.requires_stop,
		       e.cost_min_cents, e.cost_max_cents, e.time_min, e.time_max
		FROM catalog.dtc_catalog c
		LEFT JOIN catalog.dtc_vehicle_estimates e
			ON e.dtc_code = c.code AND e.model_year_id = $2
		WHERE c.code = $1
	`, code, modelYearID)

	e := &telemetrydomain.DTCCatalogWithEstimates{}
	err := row.Scan(&e.Code, &e.Description, &e.System, &e.Severity, &e.RequiresStop,
		&e.CostMinCents, &e.CostMaxCents, &e.TimeMin, &e.TimeMax)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}
