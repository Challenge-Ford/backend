package vehicleusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	"torque/internal/core/appctx"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type DeleteVehicleUseCase struct {
	repo           vehicledomain.Repository
	deviceResolver vehicledomain.DeviceResolver
}

func NewDeleteVehicle(repo vehicledomain.Repository, deviceResolver vehicledomain.DeviceResolver) *DeleteVehicleUseCase {
	return &DeleteVehicleUseCase{repo: repo, deviceResolver: deviceResolver}
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

	hasDevice, err := uc.deviceResolver.HasCommissioned(ctx, uuid.UUID(id))
	if err != nil {
		return apperr.Internal("failed to check commissioned device", err)
	}
	if hasDevice {
		return apperr.Conflict("vehicle has a commissioned device and cannot be deleted")
	}

	vehicle.Delete(auth.UserID)

	if err := uc.repo.Save(ctx, vehicle); err != nil {
		return apperr.Internal("failed to delete vehicle", err)
	}

	return nil
}
