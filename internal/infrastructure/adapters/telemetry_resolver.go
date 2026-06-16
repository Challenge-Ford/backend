package adapters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type checkActiveDTCs interface {
	Execute(ctx context.Context, vehicleIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

type TelemetryResolverAdapter struct {
	checkActiveDTCs checkActiveDTCs
}

func NewTelemetryResolver(checkActiveDTCs checkActiveDTCs) *TelemetryResolverAdapter {
	return &TelemetryResolverAdapter{checkActiveDTCs: checkActiveDTCs}
}

func (a *TelemetryResolverAdapter) HasActiveDTCs(ctx context.Context, vehicleIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	result, err := a.checkActiveDTCs.Execute(ctx, vehicleIDs)
	if err != nil {
		return nil, fmt.Errorf("telemetry resolver: %w", err)
	}
	return result, nil
}
