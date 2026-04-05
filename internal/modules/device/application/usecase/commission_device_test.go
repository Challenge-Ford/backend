package deviceusecase_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	devicedto "torque/internal/modules/device/application/dto"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicedomain "torque/internal/modules/device/domain"
	mockdevice "torque/mocks/device/domain"
)

func TestCommissionDevice_Execute(t *testing.T) {
	deviceID := devicedomain.DeviceID(uuid.New())
	vehicleID := uuid.New()
	validate := validator.New()
	validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		_, err := uuid.Parse(fl.Field().String())
		return err == nil
	})

	input := devicedto.CommissionDeviceInput{VehicleID: vehicleID.String()}

	device := &devicedomain.Device{ID: deviceID, Name: "device-01"}

	t.Run("commissions device successfully", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		resolver := mockdevice.NewMockVehicleResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(device, nil)
		resolver.EXPECT().Exists(ctx, vehicleID).Return(true, nil)
		repo.EXPECT().GetByVehicleID(ctx, vehicleID).Return(nil, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(d *devicedomain.Device) bool {
			return d.ID == deviceID && d.VehicleID != nil && *d.VehicleID == vehicleID
		})).Return(nil)

		uc := deviceusecase.NewCommissionDevice(repo, resolver, validate)
		out, err := uc.Execute(ctx, deviceID, input)

		require.NoError(t, err)
		assert.NotNil(t, out)
	})

	t.Run("returns not found when device does not exist", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		resolver := mockdevice.NewMockVehicleResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(nil, nil)

		uc := deviceusecase.NewCommissionDevice(repo, resolver, validate)
		_, err := uc.Execute(ctx, deviceID, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns not found when vehicle does not exist", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		resolver := mockdevice.NewMockVehicleResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(device, nil)
		resolver.EXPECT().Exists(ctx, vehicleID).Return(false, nil)

		uc := deviceusecase.NewCommissionDevice(repo, resolver, validate)
		_, err := uc.Execute(ctx, deviceID, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns conflict when vehicle already has a commissioned device", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		resolver := mockdevice.NewMockVehicleResolver(t)
		ctx := authCtx()

		otherDevice := &devicedomain.Device{ID: devicedomain.DeviceID(uuid.New()), Name: "other-device"}

		repo.EXPECT().GetByID(ctx, deviceID).Return(device, nil)
		resolver.EXPECT().Exists(ctx, vehicleID).Return(true, nil)
		repo.EXPECT().GetByVehicleID(ctx, vehicleID).Return(otherDevice, nil)

		uc := deviceusecase.NewCommissionDevice(repo, resolver, validate)
		_, err := uc.Execute(ctx, deviceID, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindConflict, appErr.Kind)
	})

	t.Run("returns internal error when resolver fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		resolver := mockdevice.NewMockVehicleResolver(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, deviceID).Return(device, nil)
		resolver.EXPECT().Exists(ctx, vehicleID).Return(false, assert.AnError)

		uc := deviceusecase.NewCommissionDevice(repo, resolver, validate)
		_, err := uc.Execute(ctx, deviceID, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
