package telemetryusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	mocktelemetry "torque/mocks/telemetry/domain"
)

func TestListActiveDTCs_Execute(t *testing.T) {
	ctx := context.Background()
	vehicleID := uuid.New()
	modelYearID := uuid.New()
	now := time.Now().UTC()

	t.Run("returns active DTCs enriched with catalog", func(t *testing.T) {
		repo := mocktelemetry.NewMockStateObservationRepository(t)
		catalogRepo := mocktelemetry.NewMockDTCCatalogRepository(t)
		resolver := mocktelemetry.NewMockVehicleResolver(t)
		system := "Engine"
		cat := &telemetrydomain.DTCCatalogWithEstimates{
			DTCCatalog: telemetrydomain.DTCCatalog{
				Code: "P0300", Description: "Random Misfire", System: &system,
				Severity: "high", RequiresStop: true,
			},
			CostMinCents: ptrInt(30000), CostMaxCents: ptrInt(150000),
			TimeMin: ptrInt(60), TimeMax: ptrInt(240),
		}

		resolver.EXPECT().GetModelYearIDByVehicleID(ctx, vehicleID).Return(modelYearID, nil)
		repo.EXPECT().ListActiveDTCs(ctx, vehicleID).Return([]*telemetrydomain.ActiveDTC{{Code: "P0300", Time: now}}, nil)
		catalogRepo.EXPECT().GetWithEstimates(ctx, "P0300", modelYearID).Return(cat, nil)

		uc := telemetryusecase.NewListActiveDTCs(repo, catalogRepo, resolver)
		result, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.Len(t, result.Data, 1)
		assert.Equal(t, "P0300", result.Data[0].Code)
		assert.Equal(t, "Random Misfire", result.Data[0].Description)
		assert.True(t, result.Data[0].RequiresStop)
		assert.Equal(t, ptrInt(240), result.Data[0].TimeMax)
	})
}

func ptrInt(v int) *int { return &v }
