package telemetrydomain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// StateObservationRepository persists and reads vehicle state snapshots.
type StateObservationRepository interface {
	Insert(ctx context.Context, observation *VehicleStateObservation) (bool, error)
	List(ctx context.Context, vehicleID uuid.UUID, from, to *time.Time, limit int, after *time.Time) ([]*VehicleStateObservation, error)
	Latest(ctx context.Context, vehicleID uuid.UUID) (*VehicleStateObservation, error)
	LatestPosition(ctx context.Context, vehicleID uuid.UUID) (*VehicleStateObservation, error)
	ListActiveDTCs(ctx context.Context, vehicleID uuid.UUID) ([]*ActiveDTC, error)
	HasActiveDTCs(ctx context.Context, vehicleIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

// DTCCatalogRepository is the driven port for reading diagnostic trouble code reference data.
type DTCCatalogRepository interface {
	GetWithEstimates(ctx context.Context, code string, modelYearID uuid.UUID) (*DTCCatalogWithEstimates, error)
}
