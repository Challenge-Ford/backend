package deviceusecase

import (
	"context"

	"torque/internal/core/appctx"
	"torque/internal/core/apperr"
	devicedomain "torque/internal/modules/device/domain"
)

type RevokeFunc interface {
	Revoke(ctx context.Context, serialNumber string) error
}

type DeleteDeviceUseCase struct {
	repo devicedomain.Repository
	pki  RevokeFunc
}

func NewDeleteDevice(repo devicedomain.Repository, pki RevokeFunc) *DeleteDeviceUseCase {
	return &DeleteDeviceUseCase{repo: repo, pki: pki}
}

func (uc *DeleteDeviceUseCase) Execute(ctx context.Context, id devicedomain.DeviceID) error {
	auth := appctx.MustGetAuth(ctx)

	device, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return apperr.Internal("failed to get device", err)
	}
	if device == nil {
		return apperr.NotFound("device")
	}

	if err := uc.pki.Revoke(ctx, device.CertificateSN); err != nil {
		return apperr.Internal("failed to revoke certificate", err)
	}

	device.Delete(auth.UserID)

	if err := uc.repo.Save(ctx, device); err != nil {
		return apperr.Internal("failed to delete device", err)
	}

	return nil
}
