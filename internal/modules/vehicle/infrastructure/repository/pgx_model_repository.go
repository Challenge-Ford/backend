package vehiclerepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type ModelRepository struct {
	pool *pgxpool.Pool
}

func NewModelRepository(pool *pgxpool.Pool) *ModelRepository {
	return &ModelRepository{pool: pool}
}

func (r *ModelRepository) GetModelByID(ctx context.Context, id vehicledomain.VehicleModelID) (*vehicledomain.VehicleModel, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, type
		FROM vehicle.vehicle_models
		WHERE id = $1
	`, id)

	var model vehicledomain.VehicleModel
	err := row.Scan(&model.ID, &model.Name, &model.Type)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *ModelRepository) ListModels(ctx context.Context) ([]*vehicledomain.VehicleModel, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, type
		FROM vehicle.vehicle_models
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*vehicledomain.VehicleModel
	for rows.Next() {
		m := &vehicledomain.VehicleModel{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Type); err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, rows.Err()
}

func (r *ModelRepository) GetModelYearByID(ctx context.Context, id vehicledomain.VehicleModelYearID) (*vehicledomain.VehicleModelYear, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT vmy.id, vmy.model_id, vmy.year, vmy.model_url,
		       vm.id, vm.name, vm.type
		FROM vehicle.vehicle_model_years vmy
		LEFT JOIN vehicle.vehicle_models vm ON vm.id = vmy.model_id
		WHERE vmy.id = $1
	`, id)

	var vmy vehicledomain.VehicleModelYear
	vm := &vehicledomain.VehicleModel{}
	vmy.Model = vm

	err := row.Scan(&vmy.ID, &vmy.ModelID, &vmy.Year, &vmy.ModelURL,
		&vm.ID, &vm.Name, &vm.Type)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if vm.ID == (vehicledomain.VehicleModelID{}) {
		vmy.Model = nil
	}
	return &vmy, nil
}

func (r *ModelRepository) ListModelYears(ctx context.Context, modelID vehicledomain.VehicleModelID) ([]*vehicledomain.VehicleModelYear, error) {
	// First get the years
	rows, err := r.pool.Query(ctx, `
		SELECT id, model_id, year, model_url
		FROM vehicle.vehicle_model_years
		WHERE model_id = $1
		ORDER BY year ASC
	`, modelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build index of years
	years := make([]*vehicledomain.VehicleModelYear, 0)
	yearMap := make(map[vehicledomain.VehicleModelYearID]*vehicledomain.VehicleModelYear)
	for rows.Next() {
		vmy := &vehicledomain.VehicleModelYear{}
		if err := rows.Scan(&vmy.ID, &vmy.ModelID, &vmy.Year, &vmy.ModelURL); err != nil {
			return nil, err
		}
		yearMap[vmy.ID] = vmy
		years = append(years, vmy)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(years) == 0 {
		return years, nil
	}

	// Then get the colors for all years
	ids := make([]vehicledomain.VehicleModelYearID, len(years))
	for i, y := range years {
		ids[i] = y.ID
	}

	colorRows, err := r.pool.Query(ctx, `
		SELECT id, model_year_id, name, hex
		FROM vehicle.vehicle_model_year_colors
		WHERE model_year_id = ANY($1::uuid[])
	`, ids)
	if err != nil {
		return nil, err
	}
	defer colorRows.Close()

	for colorRows.Next() {
		c := &vehicledomain.VehicleModelYearColor{}
		if err := colorRows.Scan(&c.ID, &c.ModelYearID, &c.Name, &c.Hex); err != nil {
			return nil, err
		}
		if y, ok := yearMap[c.ModelYearID]; ok {
			y.Colors = append(y.Colors, *c)
		}
	}
	return years, colorRows.Err()
}
