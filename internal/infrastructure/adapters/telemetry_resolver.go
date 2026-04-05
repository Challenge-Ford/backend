package adapters

import (
	"context"
	"fmt"
)

type checkActiveDTCs interface {
	Execute(ctx context.Context, vins []string) (map[string]bool, error)
}

type TelemetryResolverAdapter struct {
	checkActiveDTCs checkActiveDTCs
}

func NewTelemetryResolver(checkActiveDTCs checkActiveDTCs) *TelemetryResolverAdapter {
	return &TelemetryResolverAdapter{checkActiveDTCs: checkActiveDTCs}
}

func (a *TelemetryResolverAdapter) HasActiveDTCs(ctx context.Context, vins []string) (map[string]bool, error) {
	result, err := a.checkActiveDTCs.Execute(ctx, vins)
	if err != nil {
		return nil, fmt.Errorf("telemetry resolver: %w", err)
	}
	return result, nil
}
