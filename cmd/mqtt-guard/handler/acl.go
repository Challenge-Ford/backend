package handler

import (
	"encoding/json"
	"net/http"

	devicedto "torque/internal/modules/device/application/dto"
	deviceusecase "torque/internal/modules/device/application/usecase"
)

type ACLHandler struct {
	authorize *deviceusecase.AuthorizeDeviceUseCase
}

func NewACLHandler(authorize *deviceusecase.AuthorizeDeviceUseCase) *ACLHandler {
	return &ACLHandler{authorize: authorize}
}

func (h *ACLHandler) ACL(w http.ResponseWriter, r *http.Request) {
	var input devicedto.MQTTACLInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeResult(w, devicedto.MQTTDeny)
		return
	}

	result := h.authorize.Execute(r.Context(), input)
	writeResult(w, result)
}

