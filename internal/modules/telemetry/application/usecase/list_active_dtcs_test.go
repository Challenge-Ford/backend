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
	modelYearID := uuid.New()

	t.Run("returns active DTCs enriched with catalog and estimates", func(t *testing.T) {
		dtcRepo := mocktelemetry.NewMockDTCRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		now := time.Now().UTC()
		active := []*telemetrydomain.DTCEntry{
			{Time: now, DeviceID: uuid.New(), VIN: vin, Code: "P0300", Status: "opened"},
			{Time: now.Add(-time.Hour), DeviceID: uuid.New(), VIN: vin, Code: "P0420", Status: "opened"},
		}
		system := "Engine"
		catP0300 := &telemetrydomain.DTCCatalogWithEstimates{
			DTCCatalog: telemetrydomain.DTCCatalog{
				Code: "P0300", Description: "Random Misfire", System: &system,
				Severity: "high", RequiresStop: true,
			},
			CostMinCents: ptrInt(30000), CostMaxCents: ptrInt(150000),
			TimeMin: ptrInt(60), TimeMax: ptrInt(240),
		}
		catP0420 := &telemetrydomain.DTCCatalogWithEstimates{
			DTCCatalog: telemetrydomain.DTCCatalog{
				Code: "P0420", Description: "Catalyst System", System: &system,
				Severity: "medium", RequiresStop: false,
			},
			CostMinCents: ptrInt(100000), CostMaxCents: ptrInt(300000),
			TimeMin: ptrInt(120), TimeMax: ptrInt(360),
		}

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		resolver.EXPECT().GetModelYearIDByVehicleID(ctx, vehicleID).Return(modelYearID, nil)
		dtcRepo.EXPECT().ListActive(ctx, vin).Return(active, nil)
		catalogRepo.EXPECT().GetWithEstimates(ctx, "P0300", modelYearID).Return(catP0300, nil)
		catalogRepo.EXPECT().GetWithEstimates(ctx, "P0420", modelYearID).Return(catP0420, nil)

		uc := telemetryusecase.NewListActiveDTCs(dtcRepo, catalogRepo, resolver)
		result, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, "P0300", result.Data[0].Code)
		assert.Equal(t, "Random Misfire", result.Data[0].Description)
		assert.True(t, result.Data[0].RequiresStop)
		assert.Equal(t, ptrInt(30000), result.Data[0].CostMinCents)
		assert.Equal(t, ptrInt(240), result.Data[0].TimeMax)
		assert.Equal(t, "P0420", result.Data[1].Code)
		assert.Equal(t, ptrInt(100000), result.Data[1].CostMinCents)
	})

	t.Run("returns code and time when catalog has no entry", func(t *testing.T) {
		dtcRepo := mocktelemetry.NewMockDTCRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		now := time.Now().UTC()
		active := []*telemetrydomain.DTCEntry{
			{Time: now, DeviceID: uuid.New(), VIN: vin, Code: "U1234", Status: "opened"},
		}

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		resolver.EXPECT().GetModelYearIDByVehicleID(ctx, vehicleID).Return(modelYearID, nil)
		dtcRepo.EXPECT().ListActive(ctx, vin).Return(active, nil)
		catalogRepo.EXPECT().GetWithEstimates(ctx, "U1234", modelYearID).Return(nil, nil)

		uc := telemetryusecase.NewListActiveDTCs(dtcRepo, catalogRepo, resolver)
		result, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.Len(t, result.Data, 1)
		assert.Equal(t, "U1234", result.Data[0].Code)
		assert.Empty(t, result.Data[0].Description)
		assert.Nil(t, result.Data[0].CostMinCents)
	})

	t.Run("returns empty list when no active DTCs", func(t *testing.T) {
		dtcRepo := mocktelemetry.NewMockDTCRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		resolver.EXPECT().GetModelYearIDByVehicleID(ctx, vehicleID).Return(modelYearID, nil)
		dtcRepo.EXPECT().ListActive(ctx, vin).Return(nil, nil)

		uc := telemetryusecase.NewListActiveDTCs(dtcRepo, catalogRepo, resolver)
		result, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.Empty(t, result.Data)
	})

	t.Run("returns internal error when resolver fails to get VIN", func(t *testing.T) {
		dtcRepo := mocktelemetry.NewMockDTCRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return("", assert.AnError)

		uc := telemetryusecase.NewListActiveDTCs(dtcRepo, catalogRepo, resolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when resolver fails to get model year", func(t *testing.T) {
		dtcRepo := mocktelemetry.NewMockDTCRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		resolver.EXPECT().GetModelYearIDByVehicleID(ctx, vehicleID).Return(uuid.Nil, assert.AnError)

		uc := telemetryusecase.NewListActiveDTCs(dtcRepo, catalogRepo, resolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns not found when vehicle does not exist", func(t *testing.T) {
		dtcRepo := mocktelemetry.NewMockDTCRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return("", nil)

		uc := telemetryusecase.NewListActiveDTCs(dtcRepo, catalogRepo, resolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns internal error when resolver fails", func(t *testing.T) {
		dtcRepo := mocktelemetry.NewMockDTCRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return("", assert.AnError)

		uc := telemetryusecase.NewListActiveDTCs(dtcRepo, catalogRepo, resolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when ListActive fails", func(t *testing.T) {
		dtcRepo := mocktelemetry.NewMockDTCRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)

		resolver.EXPECT().GetVINByID(ctx, vehicleID).Return(vin, nil)
		resolver.EXPECT().GetModelYearIDByVehicleID(ctx, vehicleID).Return(modelYearID, nil)
		dtcRepo.EXPECT().ListActive(ctx, vin).Return(nil, assert.AnError)

		uc := telemetryusecase.NewListActiveDTCs(dtcRepo, catalogRepo, resolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}

func ptrInt(v int) *int { return &v }
