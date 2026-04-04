package telemetryusecase

import (
	"context"

	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type RecordTelemetryUseCase struct {
	repo telemetrydomain.Repository
}

func NewRecordTelemetry(repo telemetrydomain.Repository) *RecordTelemetryUseCase {
	return &RecordTelemetryUseCase{repo: repo}
}

func (uc *RecordTelemetryUseCase) Execute(ctx context.Context, entry *telemetrydomain.TelemetryEntry) error {
	return uc.repo.Insert(ctx, entry)
}
