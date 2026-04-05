package vehicledomain

import (
	"context"

	"github.com/google/uuid"
)

// DeviceResolver checks device state outside this module.
// Implemented outside this module to keep vehicle isolated.
type DeviceResolver interface {
	HasCommissioned(ctx context.Context, vehicleID uuid.UUID) (bool, error)
}

// TelemetryResolver checks telemetry state outside this module.
// Implemented outside this module to keep vehicle isolated.
type TelemetryResolver interface {
	HasActiveDTCs(ctx context.Context, vins []string) (map[string]bool, error)
}
