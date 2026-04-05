package vehicleusecase_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	vehicleusecase "torque/internal/modules/vehicle/application/usecase"
	vehicledomain "torque/internal/modules/vehicle/domain"
	mockvehicle "torque/mocks/vehicle/domain"
)

func TestListVehicleModelYears_Execute(t *testing.T) {
	modelID := uuid.New()
	model := &vehicledomain.VehicleModel{
		ID:   vehicledomain.VehicleModelID(modelID),
		Name: "Corolla",
		Type: "sedan",
	}
	years := []*vehicledomain.VehicleModelYear{
		{ID: vehicledomain.NewVehicleModelYearID(), ModelID: vehicledomain.VehicleModelID(modelID), Year: 2023},
		{ID: vehicledomain.NewVehicleModelYearID(), ModelID: vehicledomain.VehicleModelID(modelID), Year: 2024},
	}

	t.Run("lists model years successfully", func(t *testing.T) {
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().GetModelByID(ctx, vehicledomain.VehicleModelID(modelID)).Return(model, nil)
		modelRepo.EXPECT().ListModelYears(ctx, vehicledomain.VehicleModelID(modelID)).Return(years, nil)

		uc := vehicleusecase.NewListVehicleModelYears(modelRepo)
		result, err := uc.Execute(ctx, modelID)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, 2023, result.Data[0].Year)
	})

	t.Run("returns not found when model does not exist", func(t *testing.T) {
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().GetModelByID(ctx, vehicledomain.VehicleModelID(modelID)).Return(nil, nil)

		uc := vehicleusecase.NewListVehicleModelYears(modelRepo)
		_, err := uc.Execute(ctx, modelID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().GetModelByID(ctx, vehicledomain.VehicleModelID(modelID)).Return(nil, assert.AnError)

		uc := vehicleusecase.NewListVehicleModelYears(modelRepo)
		_, err := uc.Execute(ctx, modelID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
