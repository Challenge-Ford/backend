package telemetryusecase

import (
	"context"

	"torque/internal/core/apperr"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type CheckActiveDTCsUseCase struct {
	repo telemetrydomain.DTCRepository
}

func NewCheckActiveDTCs(repo telemetrydomain.DTCRepository) *CheckActiveDTCsUseCase {
	return &CheckActiveDTCsUseCase{repo: repo}
}

func (uc *CheckActiveDTCsUseCase) Execute(ctx context.Context, vins []string) (map[string]bool, error) {
	result, err := uc.repo.HasActiveDTCs(ctx, vins)
	if err != nil {
		return nil, apperr.Internal("failed to check active dtcs", err)
	}
	return result, nil
}
