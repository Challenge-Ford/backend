package messaging

import "time"

type TelemetryMessage struct {
	Time           time.Time  `json:"time"`
	VIN            string     `json:"vin"`
	Lat            *float64   `json:"lat"`
	Lng            *float64   `json:"lng"`
	Alt            *float64   `json:"alt"`
	GPSSpeed       *float64   `json:"gps_speed"`
	Heading        *float64   `json:"heading"`
	HDOP           *float64   `json:"hdop"`
	RPM            *int       `json:"rpm"`
	Speed          *int       `json:"speed"`
	CoolantTemp    *float64   `json:"coolant_temp"`
	IntakeTemp     *float64   `json:"intake_temp"`
	EngineLoad     *float64   `json:"engine_load"`
	ThrottlePos    *float64   `json:"throttle_pos"`
	FuelLevel      *float64   `json:"fuel_level"`
	FuelTrimShort  *float64   `json:"fuel_trim_short"`
	FuelTrimLong   *float64   `json:"fuel_trim_long"`
	MAF            *float64   `json:"maf"`
	BatteryVoltage *float64   `json:"battery_voltage"`
}

type DTCMessage struct {
	VIN    string    `json:"vin"`
	Code   string    `json:"code"`
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}
