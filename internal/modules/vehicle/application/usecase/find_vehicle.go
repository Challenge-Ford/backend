package vehicleusecase

import (
	"context"

	"torque/internal/core/apperr"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type FindVehicleUseCase struct {
	repo vehicledomain.Repository
}

func NewFindVehicle(repo vehicledomain.Repository) *FindVehicleUseCase {
	return &FindVehicleUseCase{repo: repo}
}

func (uc *FindVehicleUseCase) Execute(ctx context.Context, id vehicledomain.VehicleID) (*vehicledto.VehicleOutput, error) {
	vehicle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperr.Internal("failed to find vehicle", err)
	}
	if vehicle == nil {
		return nil, nil
	}
	return vehicledto.ToVehicleOutput(vehicle), nil
}
