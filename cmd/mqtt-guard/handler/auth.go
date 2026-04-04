package handler

import (
	"encoding/json"
	"net/http"

	devicedto "torque/internal/modules/device/application/dto"
	deviceusecase "torque/internal/modules/device/application/usecase"
)

type AuthHandler struct {
	authenticate *deviceusecase.AuthenticateDeviceUseCase
}

func NewAuthHandler(authenticate *deviceusecase.AuthenticateDeviceUseCase) *AuthHandler {
	return &AuthHandler{authenticate: authenticate}
}

func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var input devicedto.MQTTAuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeResult(w, devicedto.MQTTDeny)
		return
	}

	result := h.authenticate.Execute(r.Context(), input)
	writeResult(w, result)
}

func writeResult(w http.ResponseWriter, result devicedto.MQTTResult) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
