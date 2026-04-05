package adapters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type VehicleResolverAdapter struct {
	repo vehicledomain.Repository
}

func NewVehicleResolver(repo vehicledomain.Repository) *VehicleResolverAdapter {
	return &VehicleResolverAdapter{repo: repo}
}

func (a *VehicleResolverAdapter) GetVINByID(ctx context.Context, vehicleID uuid.UUID) (string, error) {
	vehicle, err := a.repo.GetByID(ctx, vehicledomain.VehicleID(vehicleID))
	if err != nil {
		return "", fmt.Errorf("vehicle resolver: %w", err)
	}
	if vehicle == nil {
		return "", nil
	}
	return string(vehicle.VIN), nil
}
