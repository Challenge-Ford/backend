package telemetrydomain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository is the driven port for persisting telemetry entries.
type Repository interface {
	Insert(ctx context.Context, entry *TelemetryEntry) error
	Latest(ctx context.Context, vin string) (*TelemetryEntry, error)
	List(ctx context.Context, vin string, from, to *time.Time, limit int, after *time.Time) ([]*TelemetryEntry, error)
	Summary(ctx context.Context, vin string, from, to time.Time, bucket string) ([]*TelemetrySummary, error)
}

// DTCRepository is the driven port for persisting DTC entries.
type DTCRepository interface {
	Insert(ctx context.Context, entry *DTCEntry) error
	ListActive(ctx context.Context, vin string) ([]*DTCEntry, error)
	HasActiveDTCs(ctx context.Context, vins []string) (map[string]bool, error)
}

// DTCCatalogRepository is the driven port for reading diagnostic trouble code reference data.
type DTCCatalogRepository interface {
	GetWithEstimates(ctx context.Context, code string, modelYearID uuid.UUID) (*DTCCatalogWithEstimates, error)
}
