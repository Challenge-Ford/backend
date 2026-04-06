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

func TestFindVehicle_Execute(t *testing.T) {
	ctx := authCtx()
	vehicle := &vehicledomain.Vehicle{
		ID:        vehicledomain.NewVehicleID(),
		VIN:       "1HGBH41JXMN109186",
		Plate:     "ABC1234",
		Color:     "#FF0000",
		ModelYear: sampleModelYear(),
	}

	t.Run("returns vehicle output when found", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		repo.EXPECT().GetByID(ctx, vehicle.ID).Return(vehicle, nil)

		uc := vehicleusecase.NewFindVehicle(repo)
		out, err := uc.Execute(ctx, vehicle.ID)

		require.NoError(t, err)
		require.NotNil(t, out)
		assert.Equal(t, vehicle.ID.String(), out.ID)
		assert.Equal(t, "1HGBH41JXMN109186", out.VIN)
		assert.Equal(t, "ABC1234", out.Plate)
		assert.Equal(t, "Corolla", out.ModelName)
		assert.Equal(t, 2024, out.Year)
	})

	t.Run("returns nil, nil when vehicle not found", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		id := vehicledomain.NewVehicleID()
		repo.EXPECT().GetByID(ctx, id).Return(nil, nil)

		uc := vehicleusecase.NewFindVehicle(repo)
		out, err := uc.Execute(ctx, id)

		require.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		id := vehicledomain.NewVehicleID()
		repo.EXPECT().GetByID(ctx, id).Return(nil, assert.AnError)

		uc := vehicleusecase.NewFindVehicle(repo)
		_, err := uc.Execute(ctx, id)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
