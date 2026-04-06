package telemetryusecase

import (
	"context"

	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type RecordDTCUseCase struct {
	repo           telemetrydomain.DTCRepository
	deviceResolver telemetrydomain.DeviceResolver
}

func NewRecordDTC(repo telemetrydomain.DTCRepository, deviceResolver telemetrydomain.DeviceResolver) *RecordDTCUseCase {
	return &RecordDTCUseCase{repo: repo, deviceResolver: deviceResolver}
}

func (uc *RecordDTCUseCase) Execute(ctx context.Context, input telemetrydto.RecordDTCInput) error {
	if input.Status != "opened" && input.Status != "closed" {
		return apperr.BadRequest("dtc status must be 'opened' or 'closed'")
	}

	device, err := uc.deviceResolver.GetCommissionedByVIN(ctx, input.VIN)
	if err != nil {
		return apperr.Internal("lookup device by vin", err)
	}
	if device == nil {
		return apperr.NotFound("commissioned device for vin " + input.VIN)
	}

	if err := uc.repo.Insert(ctx, &telemetrydomain.DTCEntry{
		Time:     input.Time.UTC(),
		DeviceID: device.ID,
		VIN:      input.VIN,
		Code:     input.Code,
		Status:   input.Status,
	}); err != nil {
		return apperr.Internal("failed to insert dtc entry", err)
	}
	return nil
}
