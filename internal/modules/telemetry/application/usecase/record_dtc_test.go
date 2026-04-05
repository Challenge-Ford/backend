package telemetryusecase_test

import (
	"context"
	"errors"
	"testing"

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

	resolvedDevice := &telemetrydomain.ResolvedDevice{ID: deviceID, VIN: vin}

	t.Run("saves open DTC when status is opened", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(resolvedDevice, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(dtc *telemetrydomain.ActiveDTC) bool {
			return dtc.Code == "P0300" && dtc.VIN == vin && !dtc.IsClosed()
		})).Return(nil)

		uc := telemetryusecase.NewRecordDTC(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordDTCInput{VIN: vin, Code: "P0300", Status: "opened"})

		require.NoError(t, err)
	})

	t.Run("saves closed DTC when status is closed", func(t *testing.T) {
		repo := mocktelemetry.NewMockDTCRepository(t)
		resolver := mocktelemetry.NewMockDeviceResolver(t)

		resolver.EXPECT().GetCommissionedByVIN(ctx, vin).Return(resolvedDevice, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(dtc *telemetrydomain.ActiveDTC) bool {
			return dtc.Code == "P0300" && dtc.IsClosed()
		})).Return(nil)

		uc := telemetryusecase.NewRecordDTC(repo, resolver)
		err := uc.Execute(ctx, telemetrydto.RecordDTCInput{VIN: vin, Code: "P0300", Status: "closed"})

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
