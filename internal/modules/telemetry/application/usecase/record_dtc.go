package telemetryusecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	devicedomain "torque/internal/modules/device/domain"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type RecordDTCUseCase struct {
	repo       telemetrydomain.DTCRepository
	deviceRepo devicedomain.Repository
}

func NewRecordDTC(repo telemetrydomain.DTCRepository, deviceRepo devicedomain.Repository) *RecordDTCUseCase {
	return &RecordDTCUseCase{repo: repo, deviceRepo: deviceRepo}
}

func (uc *RecordDTCUseCase) Execute(ctx context.Context, input telemetrydto.RecordDTCInput) error {
	if input.Status != "opened" && input.Status != "closed" {
		return apperr.BadRequest("dtc status must be 'opened' or 'closed'")
	}

	device, err := uc.deviceRepo.GetCommissionedByVIN(ctx, input.VIN)
	if err != nil {
		return apperr.Internal("lookup device by vin", err)
	}
	if device == nil {
		return apperr.NotFound("commissioned device for vin " + input.VIN)
	}

	dtc := telemetrydomain.NewActiveDTC(uuid.UUID(device.ID), input.VIN, input.Code, time.Now().UTC())
	if input.Status == "closed" {
		dtc.Close()
	}

	return uc.repo.Save(ctx, dtc)
}
