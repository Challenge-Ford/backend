package telemetrydomain

import (
	"time"

	"github.com/google/uuid"
)

type TelemetryEntry struct {
	Time           time.Time
	DeviceID       uuid.UUID
	VIN            string
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

type DTCEvent struct {
	ID       uuid.UUID
	DeviceID uuid.UUID
	VIN      string
	Code     string
	OpenedAt time.Time
	ClosedAt *time.Time
}

type Session struct {
	ID       uuid.UUID
	DeviceID uuid.UUID
	VIN      string
	OpenedAt time.Time
	ClosedAt *time.Time
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
