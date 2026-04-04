package vehicleusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	"torque/internal/core/appctx"
	devicedomain "torque/internal/modules/device/domain"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type DeleteVehicleUseCase struct {
	repo       vehicledomain.Repository
	deviceRepo devicedomain.Repository
}

func NewDeleteVehicle(repo vehicledomain.Repository, deviceRepo devicedomain.Repository) *DeleteVehicleUseCase {
	return &DeleteVehicleUseCase{repo: repo, deviceRepo: deviceRepo}
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

	device, err := uc.deviceRepo.GetByVehicleID(ctx, uuid.UUID(id))
	if err != nil {
		return apperr.Internal("failed to check commissioned device", err)
	}
	if device != nil {
		return apperr.Conflict("vehicle has a commissioned device and cannot be deleted")
	}

	vehicle.Delete(auth.UserID)

	if err := uc.repo.Save(ctx, vehicle); err != nil {
		return apperr.Internal("failed to delete vehicle", err)
	}

	return nil
}
