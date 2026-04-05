package telemetryusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	mocktelemetry "torque/mocks/telemetry/domain"
)

func TestListActiveDTCs_Execute(t *testing.T) {
	ctx := context.Background()
	vehicleID := uuid.New()
	vin := "ASD21W31231244521"

	t.Run("returns active DTCs for vehicle", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		now := time.Now().UTC()
		dtcs := []*telemetrydomain.DTCEntry{
			{Time: now, DeviceID: uuid.New(), VIN: vin, Code: "P0300", Status: "opened"},
			{Time: now.Add(-time.Hour), DeviceID: uuid.New(), VIN: vin, Code: "P0420", Status: "opened"},
		}

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		repo.EXPECT().ListActive(ctx, vin).Return(dtcs, nil)

		uc := telemetryusecase.NewListActiveDTCs(repo, resolver)
		result, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, "P0300", result.Data[0].Code)
		assert.Equal(t, "P0420", result.Data[1].Code)
	})

	t.Run("returns empty list when no active DTCs", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		repo.EXPECT().ListActive(ctx, vin).Return(nil, nil)

		uc := telemetryusecase.NewListActiveDTCs(repo, resolver)
		result, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.Empty(t, result.Data)
	})

	t.Run("returns not found when vehicle does not exist", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return("", nil)

		uc := telemetryusecase.NewListActiveDTCs(repo, resolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns internal error when resolver fails", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return("", assert.AnError)

		uc := telemetryusecase.NewListActiveDTCs(repo, resolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		repo.EXPECT().ListActive(ctx, vin).Return(nil, assert.AnError)

		uc := telemetryusecase.NewListActiveDTCs(repo, resolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
