package vehicleusecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	vehicleusecase "torque/internal/modules/vehicle/application/usecase"
	vehicledomain "torque/internal/modules/vehicle/domain"
	mockvehicle "torque/mocks/vehicle/domain"
)

func TestListVehicleModels_Execute(t *testing.T) {
	models := []*vehicledomain.VehicleModel{
		{ID: vehicledomain.NewVehicleModelID(), Name: "Corolla", Type: "sedan"},
		{ID: vehicledomain.NewVehicleModelID(), Name: "Hilux", Type: "pickup"},
	}

	t.Run("lists models successfully", func(t *testing.T) {
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().ListModels(ctx).Return(models, nil)

		uc := vehicleusecase.NewListVehicleModels(modelRepo)
		result, err := uc.Execute(ctx)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, "Corolla", result.Data[0].Name)
	})

	t.Run("returns empty list when no models", func(t *testing.T) {
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().ListModels(ctx).Return([]*vehicledomain.VehicleModel{}, nil)

		uc := vehicleusecase.NewListVehicleModels(modelRepo)
		result, err := uc.Execute(ctx)

		require.NoError(t, err)
		assert.Empty(t, result.Data)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().ListModels(ctx).Return(nil, assert.AnError)

		uc := vehicleusecase.NewListVehicleModels(modelRepo)
		_, err := uc.Execute(ctx)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
