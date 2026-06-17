package telemetryusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type GetLatestLocationUseCase struct {
	repo telemetrydomain.StateObservationRepository
}

func NewGetLatestLocation(repo telemetrydomain.StateObservationRepository) *GetLatestLocationUseCase {
	return &GetLatestLocationUseCase{repo: repo}
}

func (uc *GetLatestLocationUseCase) Execute(ctx context.Context, vehicleID uuid.UUID) (*telemetrydto.LocationOutput, error) {
	observation, err := uc.repo.LatestPosition(ctx, vehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to get latest vehicle location", err)
	}
	if observation == nil || observation.State.Position == nil {
		return nil, apperr.NotFound("vehicle location")
	}
	return telemetrydto.ToLocationOutput(observation), nil
}
