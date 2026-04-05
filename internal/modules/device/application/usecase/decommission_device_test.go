package deviceusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"torque/internal/core/appctx"
	"torque/internal/core/apperr"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicedomain "torque/internal/modules/device/domain"
	mockdevice "torque/mocks/device/domain"
)

func authCtx() context.Context {
	return appctx.WithAuth(context.Background(), appctx.AuthContext{
		UserID: uuid.New(),
		Role:   "admin",
	})
}

func TestDecommissionDevice_Execute(t *testing.T) {
	deviceID := devicedomain.DeviceID(uuid.New())
	vehicleID := uuid.New()

	t.Run("decommissions device successfully", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		ctx := authCtx()

		commissioned := &devicedomain.Device{
			ID:        deviceID,
			Name:      "device-01",
			VehicleID: &vehicleID,
		}

		repo.EXPECT().GetByID(ctx, deviceID).Return(commissioned, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(d *devicedomain.Device) bool {
			return d.ID == deviceID && d.VehicleID == nil
		})).Return(nil)

		uc := deviceusecase.NewDecommissionDevice(repo)
		out, err := uc.Execute(ctx, deviceID)

		require.NoError(t, err)
		assert.NotNil(t, out)
	})

	t.Run("returns not found when device does not exist", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(nil, nil)

		uc := deviceusecase.NewDecommissionDevice(repo)
		_, err := uc.Execute(ctx, deviceID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns conflict when device is not commissioned", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(&devicedomain.Device{
			ID:   deviceID,
			Name: "device-01",
		}, nil)

		uc := deviceusecase.NewDecommissionDevice(repo)
		_, err := uc.Execute(ctx, deviceID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindConflict, appErr.Kind)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(nil, assert.AnError)

		uc := deviceusecase.NewDecommissionDevice(repo)
		_, err := uc.Execute(ctx, deviceID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
