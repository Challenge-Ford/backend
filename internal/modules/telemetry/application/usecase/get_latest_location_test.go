package telemetryusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	mocktelemetry "torque/mocks/telemetry/domain"
)

func TestGetLatestLocation_Execute(t *testing.T) {
	ctx := context.Background()
	vehicleID := uuid.New()
	deviceID := uuid.New()
	now := time.Now().UTC()
	source := "gps"
	lat := -23.5505
	lng := -46.6333
	speed := 42.0

	t.Run("returns latest observed position", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		repo.EXPECT().LatestPosition(ctx, vehicleID).Return(&telemetrydomain.VehicleStateObservation{
			DeviceID:   deviceID,
			ObservedAt: now,
			State: telemetrydomain.VehicleState{
				Position: &telemetrydomain.PositionState{
					Source: &source,
					Lat:    &lat,
					Lng:    &lng,
					Speed:  &speed,
				},
			},
		}, nil)

		uc := telemetryusecase.NewGetLatestLocation(repo)
		out, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		require.Equal(t, now, out.Time)
		require.Equal(t, deviceID.String(), out.DeviceID)
		require.Equal(t, source, *out.Source)
		require.Equal(t, lat, *out.Lat)
		require.Equal(t, lng, *out.Lng)
		require.Equal(t, speed, *out.Speed)
	})

	t.Run("returns not found when no position is available", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		repo.EXPECT().LatestPosition(ctx, vehicleID).Return(nil, nil)

		uc := telemetryusecase.NewGetLatestLocation(repo)
		out, err := uc.Execute(ctx, vehicleID)

		require.Nil(t, out)
		require.Error(t, err)
		require.Equal(t, apperr.KindNotFound, err.(*apperr.Error).Kind)
	})
}
