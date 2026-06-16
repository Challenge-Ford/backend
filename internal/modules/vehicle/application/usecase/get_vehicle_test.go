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

func TestGetVehicle_Execute(t *testing.T) {
	vehicleID := vehicledomain.NewVehicleID()
	existing := &vehicledomain.Vehicle{
		ID:        vehicleID,
		VIN:       "1HGBH41JXMN109186",
		Plate:     "ABC1234",
		ModelYear: sampleModelYear(),
	}

	t.Run("returns vehicle with active DTC flag", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		telemetryResolver.EXPECT().HasActiveDTCs(ctx, []uuid.UUID{uuid.UUID(existing.ID)}).
			Return(map[uuid.UUID]bool{uuid.UUID(existing.ID): true}, nil)

		uc := vehicleusecase.NewGetVehicle(repo, telemetryResolver)
		out, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.Equal(t, string(existing.VIN), out.VIN)
		assert.True(t, out.HasActiveDTCs)
	})

	t.Run("returns vehicle when resolver returns empty map", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		telemetryResolver.EXPECT().HasActiveDTCs(ctx, []uuid.UUID{uuid.UUID(existing.ID)}).
			Return(map[uuid.UUID]bool{}, nil)

		uc := vehicleusecase.NewGetVehicle(repo, telemetryResolver)
		out, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.False(t, out.HasActiveDTCs)
	})

	t.Run("returns not found when vehicle does not exist", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(nil, nil)

		uc := vehicleusecase.NewGetVehicle(repo, telemetryResolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns internal error when GetByID fails", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(nil, assert.AnError)

		uc := vehicleusecase.NewGetVehicle(repo, telemetryResolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when telemetry resolver fails", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		telemetryResolver.EXPECT().HasActiveDTCs(ctx, []uuid.UUID{uuid.UUID(existing.ID)}).
			Return(nil, assert.AnError)

		uc := vehicleusecase.NewGetVehicle(repo, telemetryResolver)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
