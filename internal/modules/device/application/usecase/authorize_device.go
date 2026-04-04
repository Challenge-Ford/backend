package deviceusecase

import (
	"context"
	"fmt"

	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type AuthorizeDeviceUseCase struct {
	repo devicedomain.Repository
}

func NewAuthorizeDevice(repo devicedomain.Repository) *AuthorizeDeviceUseCase {
	return &AuthorizeDeviceUseCase{repo: repo}
}

func (uc *AuthorizeDeviceUseCase) Execute(ctx context.Context, input devicedto.MQTTACLInput) devicedto.MQTTResult {
	identity := input.Identity()
	if identity == "" || input.Topic == "" || input.Action != "publish" {
		return devicedto.MQTTDeny
	}

	device, err := uc.repo.GetByCNWithVehicle(ctx, identity)
	if err != nil || device == nil {
		return devicedto.MQTTDeny
	}

	if device.VehicleID == nil || device.VehicleVIN == nil {
		return devicedto.MQTTDeny
	}

	expected := fmt.Sprintf("torque/vehicles/%s/telemetry", *device.VehicleVIN)
	if input.Topic != expected {
		return devicedto.MQTTDeny
	}

	return devicedto.MQTTAllow
}
