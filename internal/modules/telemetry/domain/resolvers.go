package telemetrydomain

import (
	"context"

	"github.com/google/uuid"
)

// DeviceResolver resolves a commissioned device by VIN.
// Implemented outside this module to keep telemetry isolated.
type DeviceResolver interface {
	GetCommissionedByVIN(ctx context.Context, vin string) (*ResolvedDevice, error)
}

type ResolvedDevice struct {
	ID  uuid.UUID
	VIN string
}

// VehicleResolver resolves a vehicle's VIN and model year by its ID.
// Implemented outside this module to keep telemetry isolated.
type VehicleResolver interface {
	GetVINByID(ctx context.Context, vehicleID uuid.UUID) (string, error)
	GetModelYearIDByVehicleID(ctx context.Context, vehicleID uuid.UUID) (uuid.UUID, error)
}
