package deviceusecase

import (
	"context"

	"torque/internal/core/appctx"
	"torque/internal/core/apperr"
	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type DecommissionDeviceUseCase struct {
	repo devicedomain.Repository
}

func NewDecommissionDevice(repo devicedomain.Repository) *DecommissionDeviceUseCase {
	return &DecommissionDeviceUseCase{repo: repo}
}

func (uc *DecommissionDeviceUseCase) Execute(ctx context.Context, id devicedomain.DeviceID) (*devicedto.DeviceOutput, error) {
	auth := appctx.MustGetAuth(ctx)

	device, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperr.Internal("failed to get device", err)
	}
	if device == nil {
		return nil, apperr.NotFound("device")
	}
	if device.VehicleID == nil {
		return nil, apperr.Conflict("device is not commissioned")
	}

	device.VehicleID = nil
	device.UpdatedBy = auth.UserID

	if err := uc.repo.Save(ctx, device); err != nil {
		return nil, apperr.Internal("failed to decommission device", err)
	}

	return devicedto.ToDeviceOutput(device), nil
}
