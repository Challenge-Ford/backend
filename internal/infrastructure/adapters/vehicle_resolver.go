package adapters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type findVehicle interface {
	Execute(ctx context.Context, id vehicledomain.VehicleID) (*vehicledto.VehicleOutput, error)
}

type existsVehicle interface {
	Execute(ctx context.Context, id vehicledomain.VehicleID) (bool, error)
}

type VehicleResolverAdapter struct {
	findVehicle   findVehicle
	existsVehicle existsVehicle
}

func NewVehicleResolver(findVehicle findVehicle, existsVehicle existsVehicle) *VehicleResolverAdapter {
	return &VehicleResolverAdapter{findVehicle: findVehicle, existsVehicle: existsVehicle}
}

func (a *VehicleResolverAdapter) GetVINByID(ctx context.Context, vehicleID uuid.UUID) (string, error) {
	out, err := a.findVehicle.Execute(ctx, vehicledomain.VehicleID(vehicleID))
	if err != nil {
		return "", fmt.Errorf("vehicle resolver: %w", err)
	}
	if out == nil {
		return "", nil
	}
	return out.VIN, nil
}

func (a *VehicleResolverAdapter) Exists(ctx context.Context, vehicleID uuid.UUID) (bool, error) {
	exists, err := a.existsVehicle.Execute(ctx, vehicledomain.VehicleID(vehicleID))
	if err != nil {
		return false, fmt.Errorf("vehicle resolver: %w", err)
	}
	return exists, nil
}
