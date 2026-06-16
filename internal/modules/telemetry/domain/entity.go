package telemetrydomain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type VehicleStateObservation struct {
	MessageID   uuid.UUID
	DeviceID    uuid.UUID
	VehicleID   uuid.UUID
	ObservedAt  time.Time
	ReceivedAt  time.Time
	State       VehicleState
	Observation ObservationMetadata
	RawPayload  json.RawMessage
}

type VehicleState struct {
	Position    *PositionState    `json:"position,omitempty"`
	Powertrain  *PowertrainState  `json:"powertrain,omitempty"`
	Fuel        *FuelState        `json:"fuel,omitempty"`
	Electrical  *ElectricalState  `json:"electrical,omitempty"`
	Diagnostics *DiagnosticsState `json:"diagnostics,omitempty"`
}

type PositionState struct {
	Source  *string  `json:"source,omitempty"`
	Lat     *float64 `json:"lat,omitempty"`
	Lng     *float64 `json:"lng,omitempty"`
	Alt     *float64 `json:"alt,omitempty"`
	Speed   *float64 `json:"speed,omitempty"`
	Heading *float64 `json:"heading,omitempty"`
	HDOP    *float64 `json:"hdop,omitempty"`
}

type PowertrainState struct {
	RPM         *int     `json:"rpm,omitempty"`
	Speed       *int     `json:"speed,omitempty"`
	EngineLoad  *float64 `json:"engine_load,omitempty"`
	ThrottlePos *float64 `json:"throttle_pos,omitempty"`
	CoolantTemp *float64 `json:"coolant_temp,omitempty"`
	IntakeTemp  *float64 `json:"intake_temp,omitempty"`
	MAF         *float64 `json:"maf,omitempty"`
}

type FuelState struct {
	Level     *float64 `json:"level,omitempty"`
	TrimShort *float64 `json:"trim_short,omitempty"`
	TrimLong  *float64 `json:"trim_long,omitempty"`
}

type ElectricalState struct {
	BatteryVoltage *float64 `json:"battery_voltage,omitempty"`
}

type DiagnosticsState struct {
	OpenDTCs []string `json:"open_dtcs,omitempty"`
}

type ObservationMetadata struct {
	Errors []ObservationError `json:"errors,omitempty"`
}

type ObservationError struct {
	Block   string  `json:"block"`
	Code    string  `json:"code"`
	Message *string `json:"message,omitempty"`
}

type ActiveDTC struct {
	Code string
	Time time.Time
}

// DTCCatalog is a reference table with standard OBD-II diagnostic trouble codes.
type DTCCatalog struct {
	Code         string
	Description  string
	System       *string
	Severity     string
	RequiresStop bool
}

// DTCCatalogWithEstimates extends DTCCatalog with vehicle-specific cost/time estimates.
type DTCCatalogWithEstimates struct {
	DTCCatalog
	CostMinCents *int
	CostMaxCents *int
	TimeMin      *int
	TimeMax      *int
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
