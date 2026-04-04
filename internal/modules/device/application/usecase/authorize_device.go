package deviceusecase

import (
	"context"
	"fmt"
	"strings"

	devicedto "torque/internal/modules/device/application/dto"
	devicedomain "torque/internal/modules/device/domain"
)

type AuthorizeDeviceUseCase struct {
	repo       devicedomain.Repository
	serviceCNs map[string]struct{}
}

func NewAuthorizeDevice(repo devicedomain.Repository, serviceCNs []string) *AuthorizeDeviceUseCase {
	m := make(map[string]struct{}, len(serviceCNs))
	for _, cn := range serviceCNs {
		m[cn] = struct{}{}
	}
	return &AuthorizeDeviceUseCase{repo: repo, serviceCNs: m}
}

func (uc *AuthorizeDeviceUseCase) Execute(ctx context.Context, input devicedto.MQTTACLInput) devicedto.MQTTResult {
	identity := input.Identity()
	if identity == "" || input.Topic == "" {
		return devicedto.MQTTDeny
	}

	if _, ok := uc.serviceCNs[identity]; ok {
		if input.Action == "subscribe" && isServiceSubscribeTopic(input.Topic) {
			return devicedto.MQTTAllow
		}
		return devicedto.MQTTDeny
	}

	if input.Action != "publish" {
		return devicedto.MQTTDeny
	}

	device, err := uc.repo.GetByCNWithVehicle(ctx, identity)
	if err != nil || device == nil {
		return devicedto.MQTTDeny
	}

	if device.VehicleID == nil || device.VehicleVIN == nil {
		return devicedto.MQTTDeny
	}

	for _, suffix := range []string{"telemetry", "dtc", "session"} {
		if input.Topic == fmt.Sprintf("torque/vehicles/%s/%s", *device.VehicleVIN, suffix) {
			return devicedto.MQTTAllow
		}
	}

	return devicedto.MQTTDeny
}

// isServiceSubscribeTopic returns true for any vehicle topic pattern that
// services (e.g. mqtt-listener) are allowed to subscribe to.
func isServiceSubscribeTopic(topic string) bool {
	for _, suffix := range []string{"telemetry", "dtc", "session"} {
		if topic == "torque/vehicles/+/"+suffix {
			return true
		}
		if strings.HasPrefix(topic, "torque/vehicles/") && strings.HasSuffix(topic, "/"+suffix) {
			return true
		}
	}
	return false
}
