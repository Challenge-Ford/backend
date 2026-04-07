package vehiclerepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"torque/internal/core/pagination"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Save(ctx context.Context, vehicle *vehicledomain.Vehicle) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO vehicle.vehicles (id, customer_id, model_year_id, vin, plate, color, created_at, updated_at, created_by, updated_by, deleted_at, deleted_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			customer_id   = EXCLUDED.customer_id,
			model_year_id = EXCLUDED.model_year_id,
			vin           = EXCLUDED.vin,
			plate         = EXCLUDED.plate,
			color         = EXCLUDED.color,
			updated_at    = EXCLUDED.updated_at,
			updated_by    = EXCLUDED.updated_by,
			deleted_at    = EXCLUDED.deleted_at,
			deleted_by    = EXCLUDED.deleted_by
	`,
		vehicle.ID, vehicle.CustomerID, vehicle.ModelYearID,
		vehicle.VIN, vehicle.Plate, vehicle.Color,
		vehicle.CreatedAt, vehicle.UpdatedAt,
		vehicle.CreatedBy, vehicle.UpdatedBy,
		vehicle.DeletedAt, vehicle.DeletedBy,
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id vehicledomain.VehicleID) (*vehicledomain.Vehicle, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT v.id, v.customer_id, v.model_year_id, v.vin, v.plate, v.color,
		       v.created_at, v.updated_at, v.created_by, v.updated_by, v.deleted_at, v.deleted_by,
		       vmy.id, vmy.model_id, vmy.year, vmy.model_url,
		       vm.id, vm.name, vm.type
		FROM vehicle.vehicles v
		LEFT JOIN catalog.vehicle_model_years vmy ON vmy.id = v.model_year_id
		LEFT JOIN catalog.vehicle_models vm ON vm.id = vmy.model_id
		WHERE v.id = $1 AND v.deleted_at IS NULL
	`, id)

	return scanVehicle(row)
}

func (r *Repository) GetByVIN(ctx context.Context, vin vehicledomain.VIN) (*vehicledomain.Vehicle, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, customer_id, model_year_id, vin, plate, color,
		       created_at, updated_at, created_by, updated_by, deleted_at, deleted_by
		FROM vehicle.vehicles
		WHERE vin = $1 AND deleted_at IS NULL
	`, vin)

	return scanVehicleSimple(row)
}

func (r *Repository) GetByPlate(ctx context.Context, plate vehicledomain.Plate) (*vehicledomain.Vehicle, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, customer_id, model_year_id, vin, plate, color,
		       created_at, updated_at, created_by, updated_by, deleted_at, deleted_by
		FROM vehicle.vehicles
		WHERE plate = $1 AND deleted_at IS NULL
	`, plate)

	return scanVehicleSimple(row)
}

func (r *Repository) Exists(ctx context.Context, id vehicledomain.VehicleID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM vehicle.vehicles WHERE id = $1 AND deleted_at IS NULL)
	`, id).Scan(&exists)
	return exists, err
}

func (r *Repository) List(ctx context.Context, page pagination.Page) ([]*vehicledomain.Vehicle, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM vehicle.vehicles WHERE deleted_at IS NULL
	`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT v.id, v.customer_id, v.model_year_id, v.vin, v.plate, v.color,
		       v.created_at, v.updated_at, v.created_by, v.updated_by, v.deleted_at, v.deleted_by,
		       vmy.id, vmy.model_id, vmy.year, vmy.model_url,
		       vm.id, vm.name, vm.type
		FROM vehicle.vehicles v
		LEFT JOIN catalog.vehicle_model_years vmy ON vmy.id = v.model_year_id
		LEFT JOIN catalog.vehicle_models vm ON vm.id = vmy.model_id
		WHERE v.deleted_at IS NULL
		ORDER BY v.created_at DESC
		OFFSET $1 LIMIT $2
	`, page.Offset(), page.PerPage)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var vehicles []*vehicledomain.Vehicle
	for rows.Next() {
		v, err := scanVehicle(rows)
		if err != nil {
			return nil, 0, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, total, rows.Err()
}

type vehicleScanner interface {
	Scan(dest ...any) error
}

func scanVehicle(s vehicleScanner) (*vehicledomain.Vehicle, error) {
	v := &vehicledomain.Vehicle{}
	vm := &vehicledomain.VehicleModel{}
	vmy := &vehicledomain.VehicleModelYear{Model: vm}
	v.ModelYear = vmy

	err := s.Scan(
		&v.ID, &v.CustomerID, &v.ModelYearID, &v.VIN, &v.Plate, &v.Color,
		&v.CreatedAt, &v.UpdatedAt, &v.CreatedBy, &v.UpdatedBy, &v.DeletedAt, &v.DeletedBy,
		&vmy.ID, &vmy.ModelID, &vmy.Year, &vmy.ModelURL,
		&vm.ID, &vm.Name, &vm.Type,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// If no model year was joined, nil it out
	if vmy.ID == (vehicledomain.VehicleModelYearID{}) {
		v.ModelYear = nil
	} else if vm.ID == (vehicledomain.VehicleModelID{}) {
		// Year exists but no model (broken FK) — avoid returning empty Model
		vmy.Model = nil
	}
	return v, nil
}

func scanVehicleSimple(s vehicleScanner) (*vehicledomain.Vehicle, error) {
	v := &vehicledomain.Vehicle{}
	err := s.Scan(
		&v.ID, &v.CustomerID, &v.ModelYearID, &v.VIN, &v.Plate, &v.Color,
		&v.CreatedAt, &v.UpdatedAt, &v.CreatedBy, &v.UpdatedBy, &v.DeletedAt, &v.DeletedBy,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return v, nil
}
