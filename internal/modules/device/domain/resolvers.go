package devicedomain

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/pki"
)

// VehicleResolver checks vehicle existence outside this module.
// Implemented outside this module to keep device isolated.
type VehicleResolver interface {
	Exists(ctx context.Context, vehicleID uuid.UUID) (bool, error)
}

// PKI issues and revokes device TLS certificates.
// Implemented outside this module to keep device isolated.
type PKI interface {
	Issue(ctx context.Context, commonName string) (*pki.IssuedCertificate, error)
	Revoke(ctx context.Context, serialNumber string) error
}
