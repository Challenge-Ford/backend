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

func TestDeleteVehicle_Execute(t *testing.T) {
	vehicleID := vehicledomain.NewVehicleID()
	existing := &vehicledomain.Vehicle{
		ID:    vehicleID,
		VIN:   "1HGBH41JXMN109186",
		Plate: "ABC1234",
	}

	t.Run("deletes vehicle successfully", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		deviceResolver := mockvehicle.NewMockDeviceResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		deviceResolver.EXPECT().HasCommissioned(ctx, uuid.UUID(vehicleID)).Return(false, nil)
		repo.EXPECT().Save(ctx, existing).Return(nil)

		uc := vehicleusecase.NewDeleteVehicle(repo, deviceResolver)
		err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
	})

	t.Run("returns not found when vehicle does not exist", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		deviceResolver := mockvehicle.NewMockDeviceResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(nil, nil)

		uc := vehicleusecase.NewDeleteVehicle(repo, deviceResolver)
		err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns conflict when vehicle has a commissioned device", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		deviceResolver := mockvehicle.NewMockDeviceResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		deviceResolver.EXPECT().HasCommissioned(ctx, uuid.UUID(vehicleID)).Return(true, nil)

		uc := vehicleusecase.NewDeleteVehicle(repo, deviceResolver)
		err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindConflict, appErr.Kind)
	})

	t.Run("returns internal error when device resolver fails", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		deviceResolver := mockvehicle.NewMockDeviceResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		deviceResolver.EXPECT().HasCommissioned(ctx, uuid.UUID(vehicleID)).Return(false, assert.AnError)

		uc := vehicleusecase.NewDeleteVehicle(repo, deviceResolver)
		err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
