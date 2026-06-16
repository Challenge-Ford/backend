package telemetryusecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type RecordVehicleStateUseCase struct {
	repo           telemetrydomain.StateObservationRepository
	deviceResolver telemetrydomain.DeviceResolver
}

func NewRecordVehicleState(repo telemetrydomain.StateObservationRepository, deviceResolver telemetrydomain.DeviceResolver) *RecordVehicleStateUseCase {
	return &RecordVehicleStateUseCase{repo: repo, deviceResolver: deviceResolver}
}

func (uc *RecordVehicleStateUseCase) Execute(ctx context.Context, input telemetrydto.RecordVehicleStateInput) error {
	if input.SchemaVersion != 1 {
		return apperr.BadRequest("unsupported schema version")
	}
	if input.MessageID == uuid.Nil || input.DeviceID == uuid.Nil || input.VehicleID == uuid.Nil {
		return apperr.BadRequest("message_id, device_id and vehicle_id are required")
	}
	if input.ObservedAt.IsZero() {
		return apperr.BadRequest("observed_at is required")
	}
	if err := validateState(input.State); err != nil {
		return err
	}
	if err := validateObservation(input.State, input.Observation); err != nil {
		return err
	}

	ok, err := uc.deviceResolver.IsCommissionedToVehicle(ctx, input.DeviceID, input.VehicleID)
	if err != nil {
		return apperr.Internal("lookup device commissioning", err)
	}
	if !ok {
		return apperr.NotFound("commissioned device for vehicle")
	}

	if _, err := uc.repo.Insert(ctx, &telemetrydomain.VehicleStateObservation{
		MessageID:   input.MessageID,
		DeviceID:    input.DeviceID,
		VehicleID:   input.VehicleID,
		ObservedAt:  input.ObservedAt.UTC(),
		ReceivedAt:  time.Now().UTC(),
		State:       input.State,
		Observation: input.Observation,
		RawPayload:  input.RawPayload,
	}); err != nil {
		return apperr.Internal("failed to insert vehicle state observation", err)
	}
	return nil
}

func validateState(state telemetrydomain.VehicleState) error {
	if state.Position != nil {
		if !validRange(state.Position.Lat, -90, 90) {
			return apperr.BadRequest("position.lat is invalid")
		}
		if !validRange(state.Position.Lng, -180, 180) {
			return apperr.BadRequest("position.lng is invalid")
		}
		if !validMin(state.Position.Speed, 0) {
			return apperr.BadRequest("position.speed is invalid")
		}
		if !validRange(state.Position.Heading, 0, 360) {
			return apperr.BadRequest("position.heading is invalid")
		}
		if !validMin(state.Position.HDOP, 0) {
			return apperr.BadRequest("position.hdop is invalid")
		}
	}
	if state.Powertrain != nil {
		if state.Powertrain.RPM != nil && *state.Powertrain.RPM < 0 {
			return apperr.BadRequest("powertrain.rpm is invalid")
		}
		if state.Powertrain.Speed != nil && *state.Powertrain.Speed < 0 {
			return apperr.BadRequest("powertrain.speed is invalid")
		}
		if !validRange(state.Powertrain.EngineLoad, 0, 100) {
			return apperr.BadRequest("powertrain.engine_load is invalid")
		}
		if !validRange(state.Powertrain.ThrottlePos, 0, 100) {
			return apperr.BadRequest("powertrain.throttle_pos is invalid")
		}
	}
	if state.Fuel != nil && !validRange(state.Fuel.Level, 0, 100) {
		return apperr.BadRequest("fuel.level is invalid")
	}
	return nil
}

func validateObservation(state telemetrydomain.VehicleState, observation telemetrydomain.ObservationMetadata) error {
	seen := map[string]struct{}{}
	for _, e := range observation.Errors {
		if e.Block == "" || e.Code == "" {
			return apperr.BadRequest("observation.errors block and code are required")
		}
		if _, ok := seen[e.Block]; ok {
			return apperr.BadRequest("duplicate observation error block")
		}
		seen[e.Block] = struct{}{}
		if hasObservedBlock(state, e.Block) {
			return apperr.BadRequest("state block cannot also have observation error")
		}
	}
	return nil
}

func hasObservedBlock(state telemetrydomain.VehicleState, block string) bool {
	switch block {
	case "position":
		return state.Position != nil
	case "powertrain":
		return state.Powertrain != nil
	case "fuel":
		return state.Fuel != nil
	case "electrical":
		return state.Electrical != nil
	case "diagnostics":
		return state.Diagnostics != nil
	default:
		return false
	}
}

func validRange(v *float64, min, max float64) bool {
	if v == nil {
		return true
	}
	return *v >= min && *v <= max
}

func validMin(v *float64, min float64) bool {
	if v == nil {
		return true
	}
	return *v >= min
}
