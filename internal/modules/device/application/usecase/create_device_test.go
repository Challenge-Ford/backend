package deviceusecase_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	"torque/internal/core/db"
	"torque/internal/core/pki"
	devicedto "torque/internal/modules/device/application/dto"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicedomain "torque/internal/modules/device/domain"
	mockdevice "torque/mocks/device/domain"
)

func newValidate() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("device_name", func(fl validator.FieldLevel) bool {
		return devicedomain.DeviceName(fl.Field().String()).Validate() == nil
	})
	return v
}

func TestCreateDevice_Execute(t *testing.T) {
	issuedCert := &pki.IssuedCertificate{
		Certificate:  "cert-pem",
		PrivateKey:   "key-pem",
		SerialNumber: "sn-123",
	}

	input := devicedto.CreateDeviceInput{Name: "TRQ-001"}

	t.Run("creates device and issues certificate", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		repo.EXPECT().GetByName(ctx, "TRQ-001").Return(nil, nil)
		pkiMock.EXPECT().Issue(ctx, mock.AnythingOfType("string")).Return(issuedCert, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(d *devicedomain.Device) bool {
			return d.Name == "TRQ-001" && d.CertificateSN == "sn-123"
		})).Return(nil)

		uc := deviceusecase.NewCreateDevice(repo, pkiMock, newValidate())
		out, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		assert.Equal(t, "cert-pem", out.Certificate)
		assert.Equal(t, "key-pem", out.PrivateKey)
	})

	t.Run("reissues certificate for soft-deleted device with same name", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		deleted := &devicedomain.Device{
			ID:            devicedomain.NewDeviceID(),
			Name:          "TRQ-001",
			CertificateSN: "old-sn",
		}
		deleted.DeletedAt = db.SoftDeletedAt{Valid: true}

		repo.EXPECT().GetByName(ctx, "TRQ-001").Return(deleted, nil)
		pkiMock.EXPECT().Revoke(ctx, "old-sn").Return(nil)
		pkiMock.EXPECT().Issue(ctx, mock.AnythingOfType("string")).Return(issuedCert, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(d *devicedomain.Device) bool {
			return d.Name == "TRQ-001" && !d.DeletedAt.Valid
		})).Return(nil)

		uc := deviceusecase.NewCreateDevice(repo, pkiMock, newValidate())
		out, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		assert.Equal(t, "cert-pem", out.Certificate)
	})

	t.Run("returns conflict when active device with same name exists", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		repo.EXPECT().GetByName(ctx, "TRQ-001").Return(&devicedomain.Device{Name: "TRQ-001"}, nil)

		uc := deviceusecase.NewCreateDevice(repo, pkiMock, newValidate())
		_, err := uc.Execute(ctx, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindConflict, appErr.Kind)
	})

	t.Run("returns validation error for invalid device name", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		uc := deviceusecase.NewCreateDevice(repo, pkiMock, newValidate())
		_, err := uc.Execute(ctx, devicedto.CreateDeviceInput{Name: "invalid-name"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindValidation, appErr.Kind)
	})

	t.Run("returns internal error when GetByName fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		repo.EXPECT().GetByName(ctx, "TRQ-001").Return(nil, assert.AnError)

		uc := deviceusecase.NewCreateDevice(repo, pkiMock, newValidate())
		_, err := uc.Execute(ctx, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when PKI issue fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		repo.EXPECT().GetByName(ctx, "TRQ-001").Return(nil, nil)
		pkiMock.EXPECT().Issue(ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)

		uc := deviceusecase.NewCreateDevice(repo, pkiMock, newValidate())
		_, err := uc.Execute(ctx, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when Save fails", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		repo.EXPECT().GetByName(ctx, "TRQ-001").Return(nil, nil)
		pkiMock.EXPECT().Issue(ctx, mock.AnythingOfType("string")).Return(issuedCert, nil)
		repo.EXPECT().Save(ctx, mock.Anything).Return(assert.AnError)

		uc := deviceusecase.NewCreateDevice(repo, pkiMock, newValidate())
		_, err := uc.Execute(ctx, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when PKI revoke fails during reissue", func(t *testing.T) {
		repo := mockdevice.NewMockRepository(t)
		pkiMock := mockdevice.NewMockPKI(t)
		ctx := authCtx()

		deleted := &devicedomain.Device{
			ID:            devicedomain.NewDeviceID(),
			Name:          "TRQ-001",
			CertificateSN: "old-sn",
		}
		deleted.DeletedAt = db.SoftDeletedAt{Valid: true}

		repo.EXPECT().GetByName(ctx, "TRQ-001").Return(deleted, nil)
		pkiMock.EXPECT().Revoke(ctx, "old-sn").Return(assert.AnError)

		uc := deviceusecase.NewCreateDevice(repo, pkiMock, newValidate())
		_, err := uc.Execute(ctx, input)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
