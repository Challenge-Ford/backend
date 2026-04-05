package telemetrydomain

import (
	"time"

	"github.com/google/uuid"
)

type TelemetryEntry struct {
	Time           time.Time `gorm:"primaryKey;not null"`
	DeviceID       uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	VIN            string    `gorm:"not null"`
	Lat            *float64
	Lng            *float64
	Alt            *float64
	GPSSpeed       *float64
	Heading        *float64
	HDOP           *float64
	RPM            *int
	Speed          *int
	CoolantTemp    *float64
	IntakeTemp     *float64
	EngineLoad     *float64
	ThrottlePos    *float64
	FuelLevel      *float64
	FuelTrimShort  *float64
	FuelTrimLong   *float64
	MAF            *float64
	BatteryVoltage *float64
}

func (TelemetryEntry) TableName() string { return "telemetry.entries" }

// ActiveDTC represents a fault code currently present on a vehicle.
// Keyed by (device_id, code) — one row per active code per device.
type ActiveDTC struct {
	DeviceID   uuid.UUID
	VIN        string
	Code       string
	DetectedAt time.Time
	closed     bool
}

func NewActiveDTC(deviceID uuid.UUID, vin, code string, at time.Time) *ActiveDTC {
	return &ActiveDTC{DeviceID: deviceID, VIN: vin, Code: code, DetectedAt: at}
}

func (d *ActiveDTC) Close() {
	d.closed = true
}

func (d *ActiveDTC) IsClosed() bool {
	return d.closed
}

type TelemetrySummary struct {
	Bucket            time.Time
	AvgRPM            *float64
	MaxRPM            *float64
	AvgSpeed          *float64
	MaxSpeed          *float64
	AvgCoolantTemp    *float64
	MaxCoolantTemp    *float64
	AvgEngineLoad     *float64
	AvgMAF            *float64
	AvgBatteryVoltage *float64
}
