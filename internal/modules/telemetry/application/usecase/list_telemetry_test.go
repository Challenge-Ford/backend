package telemetryusecase_test

import (
	"context"
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

func TestListTelemetry_Execute(t *testing.T) {
	ctx := context.Background()
	vehicleID := uuid.New()
	vin := "ASD21W31231244521"

	t.Run("returns entries for vehicle", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		now := time.Now().UTC()
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		entries := []*telemetrydomain.TelemetryEntry{
			{Time: now, VIN: vin},
			{Time: now.Add(-time.Minute), VIN: vin},
		}

		input := telemetrydto.ListTelemetryInput{VehicleID: vehicleID, Limit: 10, From: &from, To: &to}

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		repo.EXPECT().List(ctx, vin, &from, &to, 11, input.After).Return(entries, nil)

		uc := telemetryusecase.NewListTelemetry(repo, resolver)
		result, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.Nil(t, result.Next)
	})

	t.Run("sets next cursor when results exceed limit", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		now := time.Now().UTC()
		entries := make([]*telemetrydomain.TelemetryEntry, 3)
		for i := range entries {
			entries[i] = &telemetrydomain.TelemetryEntry{Time: now.Add(-time.Duration(i) * time.Minute), VIN: vin}
		}

		input := telemetrydto.ListTelemetryInput{VehicleID: vehicleID, Limit: 2}

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		repo.EXPECT().List(ctx, vin, mock.Anything, mock.Anything, 3, input.After).Return(entries, nil)

		uc := telemetryusecase.NewListTelemetry(repo, resolver)
		result, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.NotNil(t, result.Next)
	})

	t.Run("applies default limit when limit is zero", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		input := telemetrydto.ListTelemetryInput{VehicleID: vehicleID, Limit: 0}

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		repo.EXPECT().List(ctx, vin, mock.Anything, mock.Anything, mock.MatchedBy(func(limit int) bool {
			return limit > 100
		}), input.After).Return(nil, nil)

		uc := telemetryusecase.NewListTelemetry(repo, resolver)
		result, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		assert.Empty(t, result.Data)
	})

	t.Run("returns not found when vehicle does not exist", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return("", nil)

		uc := telemetryusecase.NewListTelemetry(repo, resolver)
		_, err := uc.Execute(ctx, telemetrydto.ListTelemetryInput{VehicleID: vehicleID})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns internal error when resolver fails", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return("", assert.AnError)

		uc := telemetryusecase.NewListTelemetry(repo, resolver)
		_, err := uc.Execute(ctx, telemetrydto.ListTelemetryInput{VehicleID: vehicleID})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		input := telemetrydto.ListTelemetryInput{VehicleID: vehicleID, Limit: 10, From: &from, To: &to}

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		repo.EXPECT().List(ctx, vin, &from, &to, 11, input.After).Return(nil, assert.AnError)

		uc := telemetryusecase.NewListTelemetry(repo, resolver)
		_, err := uc.Execute(ctx, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
