package deviceusecase

import (
	"context"

	"torque/internal/core/apperr"
	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type FindCommissionedByVINUseCase struct {
	repo devicedomain.Repository
}

func NewFindCommissionedByVIN(repo devicedomain.Repository) *FindCommissionedByVINUseCase {
	return &FindCommissionedByVINUseCase{repo: repo}
}

func (uc *FindCommissionedByVINUseCase) Execute(ctx context.Context, vin string) (*devicedto.DeviceOutput, error) {
	device, err := uc.repo.GetCommissionedByVIN(ctx, vin)
	if err != nil {
		return nil, apperr.Internal("failed to find commissioned device by VIN", err)
	}
	if device == nil {
		return nil, nil
	}
	return devicedto.ToDeviceOutput(device), nil
}
