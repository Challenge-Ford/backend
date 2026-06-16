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

func TestListTelemetry_Execute(t *testing.T) {
	ctx := context.Background()
	vehicleID := uuid.New()
	now := time.Now().UTC()
	rpm := 2400
	speed := 72
	fuel := 64.0
	battery := 13.8
	entries := []*telemetrydomain.VehicleStateObservation{
		{
			MessageID:  uuid.New(),
			DeviceID:   uuid.New(),
			VehicleID:  vehicleID,
			ObservedAt: now,
			State: telemetrydomain.VehicleState{
				Powertrain: &telemetrydomain.PowertrainState{RPM: &rpm, Speed: &speed},
				Fuel:       &telemetrydomain.FuelState{Level: &fuel},
				Electrical: &telemetrydomain.ElectricalState{BatteryVoltage: &battery},
			},
		},
	}
	input := telemetrydto.ListVehicleStateInput{VehicleID: vehicleID, Limit: 10}
	repo := mocktelemetry.NewMockStateObservationRepository(t)
	repo.EXPECT().List(ctx, vehicleID, mock.Anything, mock.Anything, 11, input.After).Return(entries, nil)

	uc := telemetryusecase.NewListTelemetry(repo)
	result, err := uc.Execute(ctx, input)

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, now, result.Data[0].Time)
	assert.Equal(t, &rpm, result.Data[0].RPM)
	assert.Equal(t, &speed, result.Data[0].Speed)
	assert.Equal(t, &fuel, result.Data[0].FuelLevel)
	assert.Equal(t, &battery, result.Data[0].BatteryVoltage)
}
