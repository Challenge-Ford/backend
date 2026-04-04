package devicedomain

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/pagination"
)

type Repository interface {
	Save(ctx context.Context, device *Device) error
	List(ctx context.Context, page pagination.Page) ([]*Device, int, error)
	GetByID(ctx context.Context, id DeviceID) (*Device, error)
	GetByName(ctx context.Context, name string) (*Device, error)
	GetByCN(ctx context.Context, cn string) (*Device, error)
	GetByVehicleID(ctx context.Context, vehicleID uuid.UUID) (*Device, error)
}
