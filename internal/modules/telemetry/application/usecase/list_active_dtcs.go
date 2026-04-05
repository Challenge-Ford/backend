package telemetryusecase

import (
	"context"

	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type ListActiveDTCsUseCase struct {
	repo        telemetrydomain.DTCRepository
	vehicleRepo vehicledomain.Repository
}

func NewListActiveDTCs(repo telemetrydomain.DTCRepository, vehicleRepo vehicledomain.Repository) *ListActiveDTCsUseCase {
	return &ListActiveDTCsUseCase{repo: repo, vehicleRepo: vehicleRepo}
}

func (uc *ListActiveDTCsUseCase) Execute(ctx context.Context, vehicleID vehicledomain.VehicleID) (*telemetrydto.DTCListOutput, error) {
	vehicle, err := uc.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if vehicle == nil {
		return nil, apperr.NotFound("vehicle")
	}

	dtcs, err := uc.repo.ListActive(ctx, string(vehicle.VIN))
	if err != nil {
		return nil, apperr.Internal("failed to list active dtcs", err)
	}

	out := make([]*telemetrydto.DTCOutput, len(dtcs))
	for i, d := range dtcs {
		out[i] = &telemetrydto.DTCOutput{
			Code:       d.Code,
			DetectedAt: d.DetectedAt,
		}
	}
	return &telemetrydto.DTCListOutput{Data: out}, nil
}
