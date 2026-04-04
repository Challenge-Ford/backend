package telemetrydto

import "time"

// OBD-only output — GPS is not exposed via API.
type TelemetryOutput struct {
	Time           time.Time `json:"time"`
	RPM            *int      `json:"rpm,omitempty"`
	Speed          *int      `json:"speed,omitempty"`
	CoolantTemp    *float64  `json:"coolant_temp,omitempty"`
	IntakeTemp     *float64  `json:"intake_temp,omitempty"`
	EngineLoad     *float64  `json:"engine_load,omitempty"`
	ThrottlePos    *float64  `json:"throttle_pos,omitempty"`
	FuelLevel      *float64  `json:"fuel_level,omitempty"`
	FuelTrimShort  *float64  `json:"fuel_trim_short,omitempty"`
	FuelTrimLong   *float64  `json:"fuel_trim_long,omitempty"`
	MAF            *float64  `json:"maf,omitempty"`
	BatteryVoltage *float64  `json:"battery_voltage,omitempty"`
}

type TelemetryListOutput struct {
	Data []*TelemetryOutput `json:"data"`
	Next *time.Time         `json:"next,omitempty"`
}

type TelemetrySummaryOutput struct {
	Bucket            time.Time `json:"bucket"`
	AvgRPM            *float64  `json:"avg_rpm,omitempty"`
	MaxRPM            *float64  `json:"max_rpm,omitempty"`
	AvgSpeed          *float64  `json:"avg_speed,omitempty"`
	MaxSpeed          *float64  `json:"max_speed,omitempty"`
	AvgCoolantTemp    *float64  `json:"avg_coolant_temp,omitempty"`
	MaxCoolantTemp    *float64  `json:"max_coolant_temp,omitempty"`
	AvgEngineLoad     *float64  `json:"avg_engine_load,omitempty"`
	AvgMAF            *float64  `json:"avg_maf,omitempty"`
	AvgBatteryVoltage *float64  `json:"avg_battery_voltage,omitempty"`
}

type DTCEventOutput struct {
	ID       string     `json:"id"`
	Code     string     `json:"code"`
	OpenedAt time.Time  `json:"opened_at"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`
}

type DTCListOutput struct {
	Data []*DTCEventOutput `json:"data"`
}
