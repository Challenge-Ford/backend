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

func TestFindCommissionedByVIN_Execute(t *testing.T) {
	ctx := authCtx()
	vin := "1HGBH41JXMN109186"

	t.Run("returns device output when found", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		deviceID := devicedomain.NewDeviceID()
		vehicleUUID := uuid.New()
		deviceVIN := vin
		devicePlate := "ABC1234"
		device := &devicedomain.Device{
			ID:            deviceID,
			Name:          "TRQ-001",
			VehicleID:     &vehicleUUID,
			VehicleVIN:    &deviceVIN,
			VehiclePlate:  &devicePlate,
			CertificateCN: "TRQ-001",
			CertificateSN: "sn-123",
		}
		repo.EXPECT().GetCommissionedByVIN(ctx, vin).Return(device, nil)

		uc := deviceusecase.NewFindCommissionedByVIN(repo)
		out, err := uc.Execute(ctx, vin)

		require.NoError(t, err)
		require.NotNil(t, out)
		assert.Equal(t, deviceID.String(), out.ID)
		assert.Equal(t, "TRQ-001", out.Name)
		assert.Equal(t, "ABC1234", out.Vehicle.Plate)
	})

	t.Run("returns device with nil vehicle when not commissioned", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		device := &devicedomain.Device{
			ID:            devicedomain.NewDeviceID(),
			Name:          "TRQ-001",
			VehicleID:     nil,
			CertificateCN: "TRQ-001",
			CertificateSN: "sn-123",
		}
		repo.EXPECT().GetCommissionedByVIN(ctx, vin).Return(device, nil)

		uc := deviceusecase.NewFindCommissionedByVIN(repo)
		out, err := uc.Execute(ctx, vin)

		require.NoError(t, err)
		require.NotNil(t, out)
		assert.Nil(t, out.Vehicle)
	})

	t.Run("returns nil, nil when device not found", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		repo.EXPECT().GetCommissionedByVIN(ctx, vin).Return(nil, nil)

		uc := deviceusecase.NewFindCommissionedByVIN(repo)
		out, err := uc.Execute(ctx, vin)

		require.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("returns internal error when repository fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		repo.EXPECT().GetCommissionedByVIN(ctx, vin).Return(nil, assert.AnError)

		uc := deviceusecase.NewFindCommissionedByVIN(repo)
		_, err := uc.Execute(ctx, vin)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
