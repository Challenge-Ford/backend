package vehicledomain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, vehicle *Vehicle) error
	GetByID(ctx context.Context, id VehicleID) (*Vehicle, error)
	GetByVIN(ctx context.Context, vin VIN) (*Vehicle, error)
	ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]*Vehicle, error)
}
