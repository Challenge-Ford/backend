package telemetrydto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type ListVehicleStateInput struct {
	VehicleID uuid.UUID
	From      *time.Time
	To        *time.Time
	Limit     int
	After     *time.Time
}

type RecordVehicleStateInput struct {
	SchemaVersion int
	MessageID     uuid.UUID
	DeviceID      uuid.UUID
	VehicleID     uuid.UUID
	ObservedAt    time.Time
	State         telemetrydomain.VehicleState
	Observation   telemetrydomain.ObservationMetadata
	RawPayload    json.RawMessage
}

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

type LocationOutput struct {
	Time     time.Time `json:"time"`
	Source   *string   `json:"source,omitempty"`
	Lat      *float64  `json:"lat,omitempty"`
	Lng      *float64  `json:"lng,omitempty"`
	Alt      *float64  `json:"alt,omitempty"`
	Speed    *float64  `json:"speed,omitempty"`
	Heading  *float64  `json:"heading,omitempty"`
	HDOP     *float64  `json:"hdop,omitempty"`
	DeviceID string    `json:"deviceId"`
}

type VehicleStateOutput struct {
	MessageID   string                              `json:"messageId"`
	DeviceID    string                              `json:"deviceId"`
	VehicleID   string                              `json:"vehicleId"`
	ObservedAt  time.Time                           `json:"observedAt"`
	ReceivedAt  time.Time                           `json:"receivedAt"`
	State       telemetrydomain.VehicleState        `json:"state"`
	Observation telemetrydomain.ObservationMetadata `json:"observation"`
}

type DTCOutput struct {
	Code         string    `json:"code"`
	Time         time.Time `json:"time"`
	Description  string    `json:"description,omitempty"`
	System       *string   `json:"system,omitempty"`
	Severity     string    `json:"severity"`
	RequiresStop bool      `json:"requiresStop"`
	CostMinCents *int      `json:"costMinCents,omitempty"`
	CostMaxCents *int      `json:"costMaxCents,omitempty"`
	TimeMin      *int      `json:"timeMin,omitempty"`
	TimeMax      *int      `json:"timeMax,omitempty"`
}

type DTCListOutput struct {
	Data []*DTCOutput `json:"data"`
}

type VehicleStateListOutput struct {
	Data []*VehicleStateOutput `json:"data"`
	Next *time.Time            `json:"next,omitempty"`
}

func ToVehicleStateOutput(e *telemetrydomain.VehicleStateObservation) *VehicleStateOutput {
	return &VehicleStateOutput{
		MessageID:   e.MessageID.String(),
		DeviceID:    e.DeviceID.String(),
		VehicleID:   e.VehicleID.String(),
		ObservedAt:  e.ObservedAt,
		ReceivedAt:  e.ReceivedAt,
		State:       e.State,
		Observation: e.Observation,
	}
}

func ToTelemetryOutput(e *telemetrydomain.VehicleStateObservation) *TelemetryOutput {
	out := &TelemetryOutput{Time: e.ObservedAt}
	if e.State.Powertrain != nil {
		out.RPM = e.State.Powertrain.RPM
		out.Speed = e.State.Powertrain.Speed
		out.CoolantTemp = e.State.Powertrain.CoolantTemp
		out.IntakeTemp = e.State.Powertrain.IntakeTemp
		out.EngineLoad = e.State.Powertrain.EngineLoad
		out.ThrottlePos = e.State.Powertrain.ThrottlePos
		out.MAF = e.State.Powertrain.MAF
	}
	if e.State.Fuel != nil {
		out.FuelLevel = e.State.Fuel.Level
		out.FuelTrimShort = e.State.Fuel.TrimShort
		out.FuelTrimLong = e.State.Fuel.TrimLong
	}
	if e.State.Electrical != nil {
		out.BatteryVoltage = e.State.Electrical.BatteryVoltage
	}
	return out
}

func ToLocationOutput(e *telemetrydomain.VehicleStateObservation) *LocationOutput {
	out := &LocationOutput{
		Time:     e.ObservedAt,
		DeviceID: e.DeviceID.String(),
	}
	if e.State.Position != nil {
		out.Source = e.State.Position.Source
		out.Lat = e.State.Position.Lat
		out.Lng = e.State.Position.Lng
		out.Alt = e.State.Position.Alt
		out.Speed = e.State.Position.Speed
		out.Heading = e.State.Position.Heading
		out.HDOP = e.State.Position.HDOP
	}
	return out
}
