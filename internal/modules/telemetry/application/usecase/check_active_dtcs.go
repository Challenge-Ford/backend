package telemetryusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type CheckActiveDTCsUseCase struct {
	repo telemetrydomain.StateObservationRepository
}

func NewCheckActiveDTCs(repo telemetrydomain.StateObservationRepository) *CheckActiveDTCsUseCase {
	return &CheckActiveDTCsUseCase{repo: repo}
}

func (uc *CheckActiveDTCsUseCase) Execute(ctx context.Context, vehicleIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	result, err := uc.repo.HasActiveDTCs(ctx, vehicleIDs)
	if err != nil {
		return nil, apperr.Internal("failed to check active dtcs", err)
	}
	return result, nil
}
