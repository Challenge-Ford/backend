package vehicleusecase

import (
	"context"

	"torque/internal/core/apperr"
	"torque/internal/core/appctx"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type DeleteVehicleUseCase struct {
	repo vehicledomain.Repository
}

func NewDeleteVehicle(repo vehicledomain.Repository) *DeleteVehicleUseCase {
	return &DeleteVehicleUseCase{repo: repo}
}

func (uc *DeleteVehicleUseCase) Execute(ctx context.Context, id vehicledomain.VehicleID) error {
	auth := appctx.MustGetAuth(ctx)

	vehicle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return apperr.Internal("failed to get vehicle", err)
	}
	if vehicle == nil {
		return apperr.NotFound("vehicle")
	}

	vehicle.Delete(auth.UserID)

	if err := uc.repo.Save(ctx, vehicle); err != nil {
		return apperr.Internal("failed to delete vehicle", err)
	}

	return nil
}
