package deviceusecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"torque/internal/core/appctx"
	"torque/internal/core/apperr"
	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type CreateDeviceUseCase struct {
	repo     devicedomain.Repository
	pki      devicedomain.PKI
	validate *validator.Validate
}

func NewCreateDevice(repo devicedomain.Repository, pki devicedomain.PKI, validate *validator.Validate) *CreateDeviceUseCase {
	return &CreateDeviceUseCase{repo: repo, pki: pki, validate: validate}
}

func (uc *CreateDeviceUseCase) Execute(ctx context.Context, input devicedto.CreateDeviceInput) (*devicedto.CreateDeviceOutput, error) {
	auth := appctx.MustGetAuth(ctx)

	if err := uc.validate.Struct(input); err != nil {
		return nil, apperr.FromValidatorErr(err)
	}

	existing, err := uc.repo.GetByName(ctx, input.Name)
	if err != nil {
		return nil, apperr.Internal("failed to check device name", err)
	}
	if existing != nil && !existing.DeletedAt.Valid {
		return nil, apperr.Conflict("device with this name already exists")
	}

	var device *devicedomain.Device
	if existing != nil && existing.DeletedAt.Valid {
		if err := uc.pki.Revoke(ctx, existing.CertificateSN); err != nil {
			return nil, apperr.Internal("failed to revoke old certificate", err)
		}
		device = existing
		device.DeletedAt = gorm.DeletedAt{}
		device.DeletedBy = nil
		device.VehicleID = nil
	} else {
		id := devicedomain.NewDeviceID()
		device = &devicedomain.Device{
			ID:   id,
			Name: input.Name,
		}
		device.CreatedBy = auth.UserID
	}

	cert, err := uc.pki.Issue(ctx, device.ID.String())
	if err != nil {
		return nil, apperr.Internal("failed to issue certificate", err)
	}

	device.CertificateCN = device.ID.String()
	device.CertificateSN = cert.SerialNumber
	device.UpdatedBy = auth.UserID

	if err := uc.repo.Save(ctx, device); err != nil {
		return nil, apperr.Internal("failed to save device", err)
	}

	out := &devicedto.CreateDeviceOutput{
		DeviceOutput: *devicedto.ToDeviceOutput(device),
		Certificate:  cert.Certificate,
		PrivateKey:   cert.PrivateKey,
	}
	return out, nil
}
