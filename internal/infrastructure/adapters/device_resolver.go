package adapters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	devicedomain "torque/internal/modules/device/domain"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type DeviceResolverAdapter struct {
	repo devicedomain.Repository
}

func NewDeviceResolver(repo devicedomain.Repository) *DeviceResolverAdapter {
	return &DeviceResolverAdapter{repo: repo}
}

func (a *DeviceResolverAdapter) GetCommissionedByVIN(ctx context.Context, vin string) (*telemetrydomain.ResolvedDevice, error) {
	device, err := a.repo.GetCommissionedByVIN(ctx, vin)
	if err != nil {
		return nil, fmt.Errorf("device resolver: %w", err)
	}
	if device == nil {
		return nil, nil
	}
	return &telemetrydomain.ResolvedDevice{
		ID:  uuid.UUID(device.ID),
		VIN: vin,
	}, nil
}
