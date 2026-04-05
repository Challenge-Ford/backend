package telemetryusecase

import (
	"context"
	"time"

	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type ListTelemetryUseCase struct {
	repo        telemetrydomain.Repository
	vehicleRepo vehicledomain.Repository
}

func NewListTelemetry(repo telemetrydomain.Repository, vehicleRepo vehicledomain.Repository) *ListTelemetryUseCase {
	return &ListTelemetryUseCase{repo: repo, vehicleRepo: vehicleRepo}
}

type ListTelemetryInput struct {
	VehicleID vehicledomain.VehicleID
	From      time.Time
	To        time.Time
	Limit     int
	After     *time.Time
}

func (uc *ListTelemetryUseCase) Execute(ctx context.Context, input ListTelemetryInput) (*telemetrydto.TelemetryListOutput, error) {
	vehicle, err := uc.vehicleRepo.GetByID(ctx, input.VehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if vehicle == nil {
		return nil, apperr.NotFound("vehicle")
	}

	limit := input.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	entries, err := uc.repo.List(ctx, string(vehicle.VIN), input.From, input.To, limit+1, input.After)
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
