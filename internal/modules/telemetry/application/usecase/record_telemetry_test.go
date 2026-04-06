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
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	mocktelemetry "torque/mocks/telemetry/domain"
)

func TestRecordTelemetry_Execute(t *testing.T) {
	ctx := context.Background()
	vin := "ASD21W31231244521"
	deviceID := uuid.New()

	resolvedDevice := &telemetrydomain.ResolvedDevice{ID: deviceID, VIN: vin}

	t.Run("inserts entry when device is commissioned", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		ts := time.Now().UTC().Truncate(time.Second)
		input := telemetrydto.RecordTelemetryInput{
			Time: ts,
			VIN:  vin,
		}

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(resolvedDevice, nil)
		repo.EXPECT().Insert(ctx, &telemetrydomain.TelemetryEntry{
			Time:     ts,
			DeviceID: deviceID,
			VIN:      vin,
		}).Return(nil)

		uc := telemetryusecase.NewRecordTelemetry(repo, resolver)
		err := uc.Execute(ctx, input)

		require.NoError(t, err)
	})

	t.Run("returns not found when device is not commissioned", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(nil, nil)

		uc := telemetryusecase.NewRecordTelemetry(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordTelemetryInput{VIN: vin})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns internal error when resolver fails", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(nil, assert.AnError)

		uc := telemetryusecase.NewRecordTelemetry(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordTelemetryInput{VIN: vin})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when Insert fails", func(t *testing.T) {
		repo := mocktelemetry.NewMockRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		ts := time.Now().UTC().Truncate(time.Second)
		input := telemetrydto.RecordTelemetryInput{
			Time: ts,
			VIN:  vin,
		}

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(resolvedDevice, nil)
		repo.EXPECT().Insert(ctx, &telemetrydomain.TelemetryEntry{
			Time:     ts,
			DeviceID: deviceID,
			VIN:      vin,
		}).Return(assert.AnError)

		uc := telemetryusecase.NewRecordTelemetry(repo, resolver)
		err := uc.Execute(ctx, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
