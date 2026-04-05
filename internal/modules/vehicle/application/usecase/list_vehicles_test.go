package vehicleusecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	"torque/internal/core/pagination"
	vehicleusecase "torque/internal/modules/vehicle/application/usecase"
	vehicledomain "torque/internal/modules/vehicle/domain"
	mockvehicle "torque/mocks/vehicle/domain"
)

func TestListVehicles_Execute(t *testing.T) {
	page := pagination.Page{Page: 1, PerPage: 10}
	vehicles := []*vehicledomain.Vehicle{
		{ID: vehicledomain.NewVehicleID(), VIN: "1HGBH41JXMN109186", ModelYear: sampleModelYear()},
		{ID: vehicledomain.NewVehicleID(), VIN: "2T1BURHE0JC034761", ModelYear: sampleModelYear()},
	}
	vins := []string{string(vehicles[0].VIN), string(vehicles[1].VIN)}

	t.Run("lists vehicles with DTC flags", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().List(ctx, page).Return(vehicles, 2, nil)
		telemetryResolver.EXPECT().HasActiveDTCs(ctx, vins).
			Return(map[string]bool{string(vehicles[0].VIN): true, string(vehicles[1].VIN): false}, nil)

		uc := vehicleusecase.NewListVehicles(repo, telemetryResolver)
		result, err := uc.Execute(ctx, page)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.True(t, result.Data[0].HasActiveDTCs)
		assert.False(t, result.Data[1].HasActiveDTCs)
	})

	t.Run("returns empty list when no vehicles", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().List(ctx, page).Return([]*vehicledomain.Vehicle{}, 0, nil)
		telemetryResolver.EXPECT().HasActiveDTCs(ctx, []string{}).Return(map[string]bool{}, nil)

		uc := vehicleusecase.NewListVehicles(repo, telemetryResolver)
		result, err := uc.Execute(ctx, page)

		require.NoError(t, err)
		assert.Empty(t, result.Data)
		assert.Equal(t, 0, result.Meta.Total)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().List(ctx, page).Return(nil, 0, assert.AnError)

		uc := vehicleusecase.NewListVehicles(repo, telemetryResolver)
		_, err := uc.Execute(ctx, page)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when telemetry resolver fails", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		telemetryResolver := mockvehicle.NewMockTelemetryResolver(t)
		ctx := authCtx()

		repo.EXPECT().List(ctx, page).Return(vehicles, 2, nil)
		telemetryResolver.EXPECT().HasActiveDTCs(ctx, vins).Return(nil, assert.AnError)

		uc := vehicleusecase.NewListVehicles(repo, telemetryResolver)
		_, err := uc.Execute(ctx, page)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
