package telemetryusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type ListActiveDTCsUseCase struct {
	repo            telemetrydomain.DTCRepository
	vehicleResolver telemetrydomain.VehicleResolver
}

func NewListActiveDTCs(repo telemetrydomain.DTCRepository, vehicleResolver telemetrydomain.VehicleResolver) *ListActiveDTCsUseCase {
	return &ListActiveDTCsUseCase{repo: repo, vehicleResolver: vehicleResolver}
}

func (uc *ListActiveDTCsUseCase) Execute(ctx context.Context, vehicleID uuid.UUID) (*telemetrydto.DTCListOutput, error) {
	vin, err := uc.vehicleResolver.GetVINByID(ctx, vehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if vin == "" {
		return nil, apperr.NotFound("vehicle")
	}

	dtcs, err := uc.repo.ListActive(ctx, vin)
	if err != nil {
		return nil, apperr.Internal("failed to list active dtcs", err)
	}

	out := make([]*telemetrydto.DTCOutput, len(dtcs))
	for i, d := range dtcs {
		out[i] = &telemetrydto.DTCOutput{
			Code: d.Code,
			Time: d.Time,
		}
	}
	return &telemetrydto.DTCListOutput{Data: out}, nil
}
