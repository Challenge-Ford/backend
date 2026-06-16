package telemetryusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	mocktelemetry "torque/mocks/telemetry/domain"
)

func TestListVehicleState_Execute(t *testing.T) {
	ctx := context.Background()
	vehicleID := uuid.New()
	now := time.Now().UTC()

	t.Run("returns observations and next cursor", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		entries := []*telemetrydomain.VehicleStateObservation{
			{MessageID: uuid.New(), DeviceID: uuid.New(), VehicleID: vehicleID, ObservedAt: now},
			{MessageID: uuid.New(), DeviceID: uuid.New(), VehicleID: vehicleID, ObservedAt: now.Add(time.Minute)},
			{MessageID: uuid.New(), DeviceID: uuid.New(), VehicleID: vehicleID, ObservedAt: now.Add(2 * time.Minute)},
		}
		input := telemetrydto.ListVehicleStateInput{VehicleID: vehicleID, Limit: 2}

		repo.EXPECT().List(ctx, vehicleID, mock.Anything, mock.Anything, 3, input.After).Return(entries, nil)

		uc := telemetryusecase.NewListVehicleState(repo)
		result, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.NotNil(t, result.Next)
	})
}
