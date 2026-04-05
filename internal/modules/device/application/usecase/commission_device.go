package deviceusecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"torque/internal/core/appctx"
	"torque/internal/core/apperr"
	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type CommissionDeviceUseCase struct {
	repo            devicedomain.Repository
	vehicleResolver devicedomain.VehicleResolver
	validate        *validator.Validate
}

func NewCommissionDevice(repo devicedomain.Repository, vehicleResolver devicedomain.VehicleResolver, validate *validator.Validate) *CommissionDeviceUseCase {
	return &CommissionDeviceUseCase{repo: repo, vehicleResolver: vehicleResolver, validate: validate}
}

func (uc *CommissionDeviceUseCase) Execute(ctx context.Context, id devicedomain.DeviceID, input devicedto.CommissionDeviceInput) (*devicedto.DeviceOutput, error) {
	auth := appctx.MustGetAuth(ctx)

	if err := uc.validate.Struct(input); err != nil {
		return nil, apperr.FromValidatorErr(err)
	}

	device, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperr.Internal("failed to get device", err)
	}
	if device == nil {
		return nil, apperr.NotFound("device")
	}

	vehicleID := uuid.MustParse(input.VehicleID)

	exists, err := uc.vehicleResolver.Exists(ctx, vehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if !exists {
		return nil, apperr.NotFound("vehicle")
	}

	if existing, err := uc.repo.GetByVehicleID(ctx, vehicleID); err != nil {
		return nil, apperr.Internal("failed to check vehicle commission", err)
	} else if existing != nil && existing.ID != device.ID {
		return nil, apperr.Conflict("vehicle already has a commissioned device")
	}

	device.VehicleID = &vehicleID
	device.UpdatedBy = auth.UserID

	if err := uc.repo.Save(ctx, device); err != nil {
		return nil, apperr.Internal("failed to commission device", err)
	}

	return devicedto.ToDeviceOutput(device), nil
}
