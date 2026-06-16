package telemetrydomain

import (
	"context"

	"github.com/google/uuid"
)

// DeviceResolver checks whether a device is currently commissioned to a vehicle.
// Implemented outside this module to keep telemetry isolated.
type DeviceResolver interface {
	IsCommissionedToVehicle(ctx context.Context, deviceID, vehicleID uuid.UUID) (bool, error)
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
