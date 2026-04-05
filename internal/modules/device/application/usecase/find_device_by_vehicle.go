package deviceusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type FindDeviceByVehicleUseCase struct {
	repo devicedomain.Repository
}

func NewFindDeviceByVehicle(repo devicedomain.Repository) *FindDeviceByVehicleUseCase {
	return &FindDeviceByVehicleUseCase{repo: repo}
}

func (uc *FindDeviceByVehicleUseCase) Execute(ctx context.Context, vehicleID uuid.UUID) (*devicedto.DeviceOutput, error) {
	device, err := uc.repo.GetByVehicleID(ctx, vehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to find device by vehicle", err)
	}
	if device == nil {
		return nil, nil
	}
	return devicedto.ToDeviceOutput(device), nil
}
