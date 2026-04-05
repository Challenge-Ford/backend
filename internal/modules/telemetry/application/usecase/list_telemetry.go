package telemetryusecase

import (
	"context"
	"time"

	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type ListTelemetryUseCase struct {
	repo            telemetrydomain.Repository
	vehicleResolver telemetrydomain.VehicleResolver
}

func NewListTelemetry(repo telemetrydomain.Repository, vehicleResolver telemetrydomain.VehicleResolver) *ListTelemetryUseCase {
	return &ListTelemetryUseCase{repo: repo, vehicleResolver: vehicleResolver}
}

func (uc *ListTelemetryUseCase) Execute(ctx context.Context, input telemetrydto.ListTelemetryInput) (*telemetrydto.TelemetryListOutput, error) {
	vin, err := uc.vehicleResolver.GetVINByID(ctx, input.VehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if vin == "" {
		return nil, apperr.NotFound("vehicle")
	}

	limit := input.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	entries, err := uc.repo.List(ctx, vin, input.From, input.To, limit+1, input.After)
	if err != nil {
		return nil, apperr.Internal("failed to list telemetry", err)
	}

	var next *time.Time
	if len(entries) > limit {
		entries = entries[:limit]
		t := entries[len(entries)-1].Time
		next = &t
	}

	out := make([]*telemetrydto.TelemetryOutput, len(entries))
	for i, e := range entries {
		out[i] = telemetrydto.ToTelemetryOutput(e)
	}
	return &telemetrydto.TelemetryListOutput{Data: out, Next: next}, nil
}
