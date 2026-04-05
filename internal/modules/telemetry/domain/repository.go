package telemetrydomain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Insert(ctx context.Context, entry *TelemetryEntry) error
	Latest(ctx context.Context, vin string) (*TelemetryEntry, error)
	List(ctx context.Context, vin string, from, to time.Time, limit int, after *time.Time) ([]*TelemetryEntry, error)
	Summary(ctx context.Context, vin string, from, to time.Time, bucket string) ([]*TelemetrySummary, error)
}

type DTCRepository interface {
	SetActive(ctx context.Context, deviceID uuid.UUID, vin, code string, at time.Time) error
	SetInactive(ctx context.Context, deviceID uuid.UUID, code string) error
	ListActive(ctx context.Context, vin string) ([]*ActiveDTC, error)
}
