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

func TestFindDeviceByVehicle_Execute(t *testing.T) {
	ctx := authCtx()
	vehicleID := uuid.New()

	t.Run("returns device output when found", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		deviceName := "TRQ-001"
		deviceVIN := "1HGBH41JXMN109186"
		device := &devicedomain.Device{
			ID:            devicedomain.NewDeviceID(),
			Name:          deviceName,
			VehicleID:     &vehicleID,
			VehicleVIN:    &deviceVIN,
			CertificateCN: "TRQ-001",
			CertificateSN: "sn-123",
		}
		repo.EXPECT().GetByVehicleID(ctx, vehicleID).Return(device, nil)

		uc := deviceusecase.NewFindDeviceByVehicle(repo)
		out, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		require.NotNil(t, out)
		assert.Equal(t, deviceName, out.Name)
		assert.Equal(t, vehicleID.String(), out.Vehicle.ID)
	})

	t.Run("returns nil, nil when device not found for vehicle", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		repo.EXPECT().GetByVehicleID(ctx, vehicleID).Return(nil, nil)

		uc := deviceusecase.NewFindDeviceByVehicle(repo)
		out, err := uc.Execute(ctx, vehicleID)

		require.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		repo.EXPECT().GetByVehicleID(ctx, vehicleID).Return(nil, assert.AnError)

		uc := deviceusecase.NewFindDeviceByVehicle(repo)
		_, err := uc.Execute(ctx, vehicleID)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
