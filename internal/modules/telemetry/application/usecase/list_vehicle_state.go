package telemetryusecase

import (
	"context"
	"time"

	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type ListVehicleStateUseCase struct {
	repo telemetrydomain.StateObservationRepository
}

func NewListVehicleState(repo telemetrydomain.StateObservationRepository) *ListVehicleStateUseCase {
	return &ListVehicleStateUseCase{repo: repo}
}

func (uc *ListVehicleStateUseCase) Execute(ctx context.Context, input telemetrydto.ListVehicleStateInput) (*telemetrydto.VehicleStateListOutput, error) {
	limit := input.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	entries, err := uc.repo.List(ctx, input.VehicleID, input.From, input.To, limit+1, input.After)
	if err != nil {
		return nil, apperr.Internal("failed to list vehicle state observations", err)
	}

	var next *time.Time
	if len(entries) > limit {
		entries = entries[:limit]
		t := entries[len(entries)-1].ObservedAt
		next = &t
	}

	out := make([]*telemetrydto.VehicleStateOutput, len(entries))
	for i, e := range entries {
		out[i] = telemetrydto.ToVehicleStateOutput(e)
	}
	return &telemetrydto.VehicleStateListOutput{Data: out, Next: next}, nil
}
