package deviceusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	"torque/internal/core/pagination"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicedomain "torque/internal/modules/device/domain"
	mockdevice "torque/mocks/device/domain"
)

func TestListDevices_Execute(t *testing.T) {
	ctx := context.Background()

	devices := []*devicedomain.Device{
		{ID: devicedomain.DeviceID(uuid.New()), Name: "TRQ-001"},
		{ID: devicedomain.DeviceID(uuid.New()), Name: "TRQ-002"},
	}

	t.Run("returns paginated devices", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)

		page := pagination.Page{Page: 1, PerPage: 10}
		repo.EXPECT().List(ctx, page).Return(devices, 2, nil)

		uc := deviceusecase.NewListDevices(repo)
		result, err := uc.Execute(ctx, page)

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, 2, result.Meta.Total)
	})

	t.Run("returns empty list when no devices exist", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)

		page := pagination.Page{Page: 1, PerPage: 10}
		repo.EXPECT().List(ctx, page).Return(nil, 0, nil)

		uc := deviceusecase.NewListDevices(repo)
		result, err := uc.Execute(ctx, page)

		require.NoError(t, err)
		assert.Empty(t, result.Data)
		assert.Equal(t, 0, result.Meta.Total)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)

		page := pagination.Page{Page: 1, PerPage: 10}
		repo.EXPECT().List(ctx, page).Return(nil, 0, assert.AnError)

		uc := deviceusecase.NewListDevices(repo)
		_, err := uc.Execute(ctx, page)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
