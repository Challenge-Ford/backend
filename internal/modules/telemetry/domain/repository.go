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
	ActiveByDevice(ctx context.Context, deviceID uuid.UUID) ([]*DTCEvent, error)
	ActiveByVIN(ctx context.Context, vin string) ([]*DTCEvent, error)
	OpenMany(ctx context.Context, deviceID uuid.UUID, vin string, codes []string, at time.Time) error
	CloseMany(ctx context.Context, deviceID uuid.UUID, codes []string, at time.Time) error
	CloseAllByDevice(ctx context.Context, deviceID uuid.UUID, at time.Time) error
	List(ctx context.Context, vin string, from, to time.Time, activeOnly bool) ([]*DTCEvent, error)
}

type SessionRepository interface {
	Open(ctx context.Context, session *Session) error
	CloseActive(ctx context.Context, deviceID uuid.UUID, at time.Time) error
	ActiveByDevice(ctx context.Context, deviceID uuid.UUID) (*Session, error)
}
