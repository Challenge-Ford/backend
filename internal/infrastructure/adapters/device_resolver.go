package adapters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	devicedto "torque/internal/modules/device/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type findCommissionedByVIN interface {
	Execute(ctx context.Context, vin string) (*devicedto.DeviceOutput, error)
}

type findDeviceByVehicle interface {
	Execute(ctx context.Context, vehicleID uuid.UUID) (*devicedto.DeviceOutput, error)
}

type DeviceResolverAdapter struct {
	findCommissionedByVIN findCommissionedByVIN
	findDeviceByVehicle   findDeviceByVehicle
}

func NewDeviceResolver(findCommissionedByVIN findCommissionedByVIN, findDeviceByVehicle findDeviceByVehicle) *DeviceResolverAdapter {
	return &DeviceResolverAdapter{
		findCommissionedByVIN: findCommissionedByVIN,
		findDeviceByVehicle:   findDeviceByVehicle,
	}
}

func (a *DeviceResolverAdapter) GetCommissionedByVIN(ctx context.Context, vin string) (*telemetrydomain.ResolvedDevice, error) {
	out, err := a.findCommissionedByVIN.Execute(ctx, vin)
	if err != nil {
		return nil, fmt.Errorf("device resolver: %w", err)
	}
	if out == nil {
		return nil, nil
	}
	id, err := uuid.Parse(out.ID)
	if err != nil {
		return nil, fmt.Errorf("device resolver: invalid device id: %w", err)
	}
	return &telemetrydomain.ResolvedDevice{ID: id, VIN: vin}, nil
}

func (a *DeviceResolverAdapter) HasCommissioned(ctx context.Context, vehicleID uuid.UUID) (bool, error) {
	out, err := a.findDeviceByVehicle.Execute(ctx, vehicleID)
	if err != nil {
		return false, fmt.Errorf("device resolver: %w", err)
	}
	return out != nil, nil
}
