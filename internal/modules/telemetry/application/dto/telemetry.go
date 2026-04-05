package telemetrydto

import (
	"time"

	"github.com/google/uuid"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type ListTelemetryInput struct {
	VehicleID uuid.UUID
	From      time.Time
	To        time.Time
	Limit     int
	After     *time.Time
}

type RecordTelemetryInput struct {
	Time           time.Time
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

type RecordDTCInput struct {
	VIN    string
	Code   string
	Status string // "opened" or "closed"
	Time   time.Time
}

// OBD-only output — GPS is not exposed via API.
type TelemetryOutput struct {
	Time           time.Time `json:"time"`
	RPM            *int      `json:"rpm,omitempty"`
	Speed          *int      `json:"speed,omitempty"`
	CoolantTemp    *float64  `json:"coolantTemp,omitempty"`
	IntakeTemp     *float64  `json:"intakeTemp,omitempty"`
	EngineLoad     *float64  `json:"engineLoad,omitempty"`
	ThrottlePos    *float64  `json:"throttlePos,omitempty"`
	FuelLevel      *float64  `json:"fuelLevel,omitempty"`
	FuelTrimShort  *float64  `json:"fuelTrimShort,omitempty"`
	FuelTrimLong   *float64  `json:"fuelTrimLong,omitempty"`
	MAF            *float64  `json:"maf,omitempty"`
	BatteryVoltage *float64  `json:"batteryVoltage,omitempty"`
}

type TelemetryListOutput struct {
	Data []*TelemetryOutput `json:"data"`
	Next *time.Time         `json:"next,omitempty"`
}

type TelemetrySummaryOutput struct {
	Bucket            time.Time `json:"bucket"`
	AvgRPM            *float64  `json:"avgRpm,omitempty"`
	MaxRPM            *float64  `json:"maxRpm,omitempty"`
	AvgSpeed          *float64  `json:"avgSpeed,omitempty"`
	MaxSpeed          *float64  `json:"maxSpeed,omitempty"`
	AvgCoolantTemp    *float64  `json:"avgCoolantTemp,omitempty"`
	MaxCoolantTemp    *float64  `json:"maxCoolantTemp,omitempty"`
	AvgEngineLoad     *float64  `json:"avgEngineLoad,omitempty"`
	AvgMAF            *float64  `json:"avgMaf,omitempty"`
	AvgBatteryVoltage *float64  `json:"avgBatteryVoltage,omitempty"`
}

type DTCOutput struct {
	Code       string    `json:"code"`
	DetectedAt time.Time `json:"detectedAt"`
}

type DTCListOutput struct {
	Data []*DTCOutput `json:"data"`
}

func ToTelemetryOutput(e *telemetrydomain.TelemetryEntry) *TelemetryOutput {
	return &TelemetryOutput{
		Time:           e.Time,
		RPM:            e.RPM,
		Speed:          e.Speed,
		CoolantTemp:    e.CoolantTemp,
		IntakeTemp:     e.IntakeTemp,
		EngineLoad:     e.EngineLoad,
		ThrottlePos:    e.ThrottlePos,
		FuelLevel:      e.FuelLevel,
		FuelTrimShort:  e.FuelTrimShort,
		FuelTrimLong:   e.FuelTrimLong,
		MAF:            e.MAF,
		BatteryVoltage: e.BatteryVoltage,
	}
}
