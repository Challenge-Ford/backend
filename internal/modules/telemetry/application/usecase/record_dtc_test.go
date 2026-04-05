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

func TestRecordDTC_Execute(t *testing.T) {
	ctx := context.Background()
	vin := "ASD21W31231244521"
	deviceID := uuid.New()
	ts := time.Now().UTC().Truncate(time.Second)

	resolvedDevice := &telemetrydomain.ResolvedDevice{ID: deviceID, VIN: vin}

	t.Run("inserts opened entry", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(resolvedDevice, nil)
		repo.EXPECT().Insert(ctx, mock.MatchedBy(func(e *telemetrydomain.DTCEntry) bool {
			return e.Code == "P0300" && e.Status == "opened" && e.VIN == vin && e.DeviceID == deviceID
		})).Return(nil)

		uc := telemetryusecase.NewRecordDTC(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordDTCInput{VIN: vin, Code: "P0300", Status: "opened", Time: ts})

		require.NoError(t, err)
	})

	t.Run("inserts closed entry", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(resolvedDevice, nil)
		repo.EXPECT().Insert(ctx, mock.MatchedBy(func(e *telemetrydomain.DTCEntry) bool {
			return e.Code == "P0300" && e.Status == "closed"
		})).Return(nil)

		uc := telemetryusecase.NewRecordDTC(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordDTCInput{VIN: vin, Code: "P0300", Status: "closed", Time: ts})

		require.NoError(t, err)
	})

	t.Run("returns bad request when status is invalid", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		uc := telemetryusecase.NewRecordDTC(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordDTCInput{VIN: vin, Code: "P0300", Status: "unknown"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindBadRequest, appErr.Kind)
	})

	t.Run("returns not found when device is not commissioned", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(nil, nil)

		uc := telemetryusecase.NewRecordDTC(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordDTCInput{VIN: vin, Code: "P0300", Status: "opened"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns internal error when resolver fails", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(nil, assert.AnError)

		uc := telemetryusecase.NewRecordDTC(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordDTCInput{VIN: vin, Code: "P0300", Status: "opened"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
