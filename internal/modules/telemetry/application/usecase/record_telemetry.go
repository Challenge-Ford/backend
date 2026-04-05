package telemetryusecase

import (
	"context"

	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type RecordTelemetryUseCase struct {
	repo           telemetrydomain.Repository
	deviceResolver telemetrydomain.DeviceResolver
}

func NewRecordTelemetry(repo telemetrydomain.Repository, deviceResolver telemetrydomain.DeviceResolver) *RecordTelemetryUseCase {
	return &RecordTelemetryUseCase{repo: repo, deviceResolver: deviceResolver}
}

func (uc *RecordTelemetryUseCase) Execute(ctx context.Context, input telemetrydto.RecordTelemetryInput) error {
	device, err := uc.deviceResolver.GetCommissionedByVIN(ctx, input.VIN)
	if err != nil {
		return apperr.Internal("lookup device by vin", err)
	}
	if device == nil {
		return apperr.NotFound("commissioned device for vin " + input.VIN)
	}

	return uc.repo.Insert(ctx, &telemetrydomain.TelemetryEntry{
		Time:           input.Time.UTC(),
		DeviceID:       device.ID,
		VIN:            input.VIN,
		Lat:            input.Lat,
		Lng:            input.Lng,
		Alt:            input.Alt,
		GPSSpeed:       input.GPSSpeed,
		Heading:        input.Heading,
		HDOP:           input.HDOP,
		RPM:            input.RPM,
		Speed:          input.Speed,
		CoolantTemp:    input.CoolantTemp,
		IntakeTemp:     input.IntakeTemp,
		EngineLoad:     input.EngineLoad,
		ThrottlePos:    input.ThrottlePos,
		FuelLevel:      input.FuelLevel,
		FuelTrimShort:  input.FuelTrimShort,
		FuelTrimLong:   input.FuelTrimLong,
		MAF:            input.MAF,
		BatteryVoltage: input.BatteryVoltage,
	})
}
