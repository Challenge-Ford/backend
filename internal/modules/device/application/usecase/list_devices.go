package deviceusecase

import (
	"context"

	"torque/internal/core/apperr"
	"torque/internal/core/pagination"
	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type ListDevicesUseCase struct {
	repo devicedomain.Repository
}

func NewListDevices(repo devicedomain.Repository) *ListDevicesUseCase {
	return &ListDevicesUseCase{repo: repo}
}

func (uc *ListDevicesUseCase) Execute(ctx context.Context, page pagination.Page) (*pagination.Result[*devicedto.DeviceOutput], error) {
	page.Normalize(pagination.DefaultConfig)

	devices, total, err := uc.repo.List(ctx, page)
	if err != nil {
		return nil, apperr.Internal("failed to list devices", err)
	}

	output := make([]*devicedto.DeviceOutput, len(devices))
	for i, d := range devices {
		output[i] = devicedto.ToDeviceOutput(d)
	}

	result := pagination.NewResult(output, page, total)
	return &result, nil
}
