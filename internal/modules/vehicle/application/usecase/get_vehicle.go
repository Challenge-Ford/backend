package vehicleusecase

import (
	"context"

	"torque/internal/core/apperr"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type GetVehicleUseCase struct {
	repo vehicledomain.Repository
}

func NewGetVehicle(repo vehicledomain.Repository) *GetVehicleUseCase {
	return &GetVehicleUseCase{repo: repo}
}

func (uc *GetVehicleUseCase) Execute(ctx context.Context, id vehicledomain.VehicleID) (*vehicledto.VehicleOutput, error) {
	vehicle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if vehicle == nil {
		return nil, apperr.NotFound("vehicle")
	}

	return vehicledto.ToVehicleOutput(vehicle), nil
}
