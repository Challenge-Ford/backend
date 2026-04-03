package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"torque/cmd/api/httperr"
	vehicleusecase "torque/internal/modules/vehicle/application/usecase"
)

type VehicleModelHandler struct {
	list      *vehicleusecase.ListVehicleModelsUseCase
	listYears *vehicleusecase.ListVehicleModelYearsUseCase
}

func NewVehicleModelHandler(
	list *vehicleusecase.ListVehicleModelsUseCase,
	listYears *vehicleusecase.ListVehicleModelYearsUseCase,
) *VehicleModelHandler {
	return &VehicleModelHandler{list: list, listYears: listYears}
}

func (h *VehicleModelHandler) List(w http.ResponseWriter, r *http.Request) {
	output, err := h.list.Execute(r.Context())
	if err != nil {
		httperr.Write(w, err)
		return
	}
	httperr.JSON(w, http.StatusOK, output)
}

func (h *VehicleModelHandler) ListYears(w http.ResponseWriter, r *http.Request) {
	modelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, err)
		return
	}
	output, err := h.listYears.Execute(r.Context(), modelID)
	if err != nil {
		httperr.Write(w, err)
		return
	}
	httperr.JSON(w, http.StatusOK, output)
}
