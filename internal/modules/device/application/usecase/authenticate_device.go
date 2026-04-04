package deviceusecase

import (
	"context"

	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type AuthenticateDeviceUseCase struct {
	repo       devicedomain.Repository
	serviceCNs map[string]struct{}
}

func NewAuthenticateDevice(repo devicedomain.Repository, serviceCNs []string) *AuthenticateDeviceUseCase {
	m := make(map[string]struct{}, len(serviceCNs))
	for _, cn := range serviceCNs {
		m[cn] = struct{}{}
	}
	return &AuthenticateDeviceUseCase{repo: repo, serviceCNs: m}
}

func (uc *AuthenticateDeviceUseCase) Execute(ctx context.Context, input devicedto.MQTTAuthInput) devicedto.MQTTResult {
	identity := input.Identity()
	if identity == "" {
		return devicedto.MQTTDeny
	}

	if _, ok := uc.serviceCNs[identity]; ok {
		return devicedto.MQTTAllow
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
