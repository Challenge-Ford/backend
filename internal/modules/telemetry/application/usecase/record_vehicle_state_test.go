package telemetryusecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	mocktelemetry "torque/mocks/telemetry/domain"
)

func TestRecordVehicleState_Execute(t *testing.T) {
	ctx := context.Background()
	deviceID := uuid.New()
	vehicleID := uuid.New()
	messageID := uuid.New()
	observedAt := time.Now().UTC().Truncate(time.Second)

	t.Run("inserts valid observation", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)
		rpm := 2400
		speed := 72
		input := telemetrydto.RecordVehicleStateInput{
			SchemaVersion: 1,
			MessageID:     messageID,
			DeviceID:      deviceID,
			VehicleID:     vehicleID,
			ObservedAt:    observedAt,
			State: telemetrydomain.VehicleState{
				Powertrain: &telemetrydomain.PowertrainState{RPM: &rpm, Speed: &speed},
			},
		}

		resolver.EXPECT().IsCommissionedToVehicle(ctx, deviceID, vehicleID).Return(true, nil)
		repo.EXPECT().Insert(ctx, mock.MatchedBy(func(e *telemetrydomain.VehicleStateObservation) bool {
			return e.MessageID == messageID && e.DeviceID == deviceID && e.VehicleID == vehicleID && e.ObservedAt.Equal(observedAt)
		})).Return(true, nil)

		uc := telemetryusecase.NewRecordVehicleState(repo, resolver)
		err := uc.Execute(ctx, input)

		require.NoError(t, err)
	})

	t.Run("rejects device not commissioned to vehicle", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().IsCommissionedToVehicle(ctx, deviceID, vehicleID).Return(false, nil)

		uc := telemetryusecase.NewRecordVehicleState(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordVehicleStateInput{
			SchemaVersion: 1,
			MessageID:     messageID,
			DeviceID:      deviceID,
			VehicleID:     vehicleID,
			ObservedAt:    observedAt,
		})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("rejects success and error for same block", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)
		lat := -23.55

		uc := telemetryusecase.NewRecordVehicleState(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordVehicleStateInput{
			SchemaVersion: 1,
			MessageID:     messageID,
			DeviceID:      deviceID,
			VehicleID:     vehicleID,
			ObservedAt:    observedAt,
			State: telemetrydomain.VehicleState{
				Position: &telemetrydomain.PositionState{Lat: &lat},
			},
			Observation: telemetrydomain.ObservationMetadata{
				Errors: []telemetrydomain.ObservationError{{Block: "position", Code: "gps_no_fix"}},
			},
		})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindBadRequest, appErr.Kind)
	})

	t.Run("rejects invalid range", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)
		lat := 200.0

		uc := telemetryusecase.NewRecordVehicleState(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordVehicleStateInput{
			SchemaVersion: 1,
			MessageID:     messageID,
			DeviceID:      deviceID,
			VehicleID:     vehicleID,
			ObservedAt:    observedAt,
			State: telemetrydomain.VehicleState{
				Position: &telemetrydomain.PositionState{Lat: &lat},
			},
		})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindBadRequest, appErr.Kind)
	})

	t.Run("preserves observed empty diagnostics", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)
		input := telemetrydto.RecordVehicleStateInput{
			SchemaVersion: 1,
			MessageID:     messageID,
			DeviceID:      deviceID,
			VehicleID:     vehicleID,
			ObservedAt:    observedAt,
			State: telemetrydomain.VehicleState{
				Diagnostics: &telemetrydomain.DiagnosticsState{OpenDTCs: []string{}},
			},
		}

		resolver.EXPECT().IsCommissionedToVehicle(ctx, deviceID, vehicleID).Return(true, nil)
		repo.EXPECT().Insert(ctx, mock.MatchedBy(func(e *telemetrydomain.VehicleStateObservation) bool {
			raw, err := json.Marshal(e.State)
			return err == nil && string(raw) == `{"diagnostics":{"open_dtcs":[]}}`
		})).Return(true, nil)

		uc := telemetryusecase.NewRecordVehicleState(repo, resolver)
		err := uc.Execute(ctx, input)

		require.NoError(t, err)
	})
}
