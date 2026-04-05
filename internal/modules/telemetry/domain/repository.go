package telemetrydomain

import (
	"context"
	"time"
)

// Repository is the driven port for persisting telemetry entries.
type Repository interface {
	Insert(ctx context.Context, entry *TelemetryEntry) error
	Latest(ctx context.Context, vin string) (*TelemetryEntry, error)
	List(ctx context.Context, vin string, from, to time.Time, limit int, after *time.Time) ([]*TelemetryEntry, error)
	Summary(ctx context.Context, vin string, from, to time.Time, bucket string) ([]*TelemetrySummary, error)
}

// DTCRepository is the driven port for persisting DTC records.
type DTCRepository interface {
	Save(ctx context.Context, dtc *ActiveDTC) error
	ListActive(ctx context.Context, vin string) ([]*ActiveDTC, error)
	HasActiveDTCs(ctx context.Context, vins []string) (map[string]bool, error)
}
