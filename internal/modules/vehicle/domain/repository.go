package vehicledomain

import (
	"context"

	"torque/internal/core/pagination"
)

type Repository interface {
	Save(ctx context.Context, vehicle *Vehicle) error
	GetByID(ctx context.Context, id VehicleID) (*Vehicle, error)
	GetByVIN(ctx context.Context, vin VIN) (*Vehicle, error)
	GetByPlate(ctx context.Context, plate Plate) (*Vehicle, error)
	List(ctx context.Context, page pagination.Page) ([]*Vehicle, int, error)
	Exists(ctx context.Context, id VehicleID) (bool, error)
}
