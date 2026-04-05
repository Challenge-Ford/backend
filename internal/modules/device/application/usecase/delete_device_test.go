package deviceusecase_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicedomain "torque/internal/modules/device/domain"
	mockdevice "torque/mocks/device/domain"
)

func TestDeleteDevice_Execute(t *testing.T) {
	deviceID := devicedomain.DeviceID(uuid.New())

	device := &devicedomain.Device{
		ID:            deviceID,
		Name:          "TRQ-001",
		CertificateSN: "sn-123",
	}

	t.Run("deletes device and revokes certificate", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(device, nil)
		pkiMock.EXPECT().Revoke(ctx, "sn-123").Return(nil)
		repo.EXPECT().Save(ctx, device).Return(nil)

		uc := deviceusecase.NewDeleteDevice(repo, pkiMock)
		err := uc.Execute(ctx, deviceID)

		require.NoError(t, err)
		assert.True(t, device.DeletedAt.Valid)
	})

	t.Run("returns not found when device does not exist", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(nil, nil)

		uc := deviceusecase.NewDeleteDevice(repo, pkiMock)
		err := uc.Execute(ctx, deviceID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns internal error when revoke fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(device, nil)
		pkiMock.EXPECT().Revoke(ctx, "sn-123").Return(assert.AnError)

		uc := deviceusecase.NewDeleteDevice(repo, pkiMock)
		err := uc.Execute(ctx, deviceID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
