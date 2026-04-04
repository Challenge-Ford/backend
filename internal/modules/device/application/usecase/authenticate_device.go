package deviceusecase

import (
	"context"

	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type AuthenticateDeviceUseCase struct {
	repo devicedomain.Repository
}

func NewAuthenticateDevice(repo devicedomain.Repository) *AuthenticateDeviceUseCase {
	return &AuthenticateDeviceUseCase{repo: repo}
}

func (uc *AuthenticateDeviceUseCase) Execute(ctx context.Context, input devicedto.MQTTAuthInput) devicedto.MQTTResult {
	identity := input.Identity()
	if identity == "" {
		return devicedto.MQTTDeny
	}

	device, err := uc.repo.GetByCN(ctx, identity)
	if err != nil || device == nil {
		return devicedto.MQTTDeny
	}

	if device.VehicleID == nil {
		return devicedto.MQTTDeny
	}

	return devicedto.MQTTAllow
}
