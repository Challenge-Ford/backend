package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"torque/cmd/api/httperr"
	"torque/internal/core/pagination"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicleusecase "torque/internal/modules/vehicle/application/usecase"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type VehicleHandler struct {
	create *vehicleusecase.CreateVehicleUseCase
	get    *vehicleusecase.GetVehicleUseCase
	list   *vehicleusecase.ListVehiclesUseCase
	update *vehicleusecase.UpdateVehicleUseCase
	delete *vehicleusecase.DeleteVehicleUseCase
}

func NewVehicleHandler(
	create *vehicleusecase.CreateVehicleUseCase,
	get *vehicleusecase.GetVehicleUseCase,
	list *vehicleusecase.ListVehiclesUseCase,
	update *vehicleusecase.UpdateVehicleUseCase,
	delete *vehicleusecase.DeleteVehicleUseCase,
) *VehicleHandler {
	return &VehicleHandler{create: create, get: get, list: list, update: update, delete: delete}
}

func (h *VehicleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input vehicledto.CreateVehicleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httperr.Write(w, err)
		return
	}

	output, err := h.create.Execute(r.Context(), input)
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusCreated, output)
}

func (h *VehicleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	output, err := h.get.Execute(r.Context(), vehicledomain.VehicleID(id))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, output)
}

func (h *VehicleHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))

	result, err := h.list.Execute(r.Context(), pagination.Page{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, result)
}

func (h *VehicleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	var input vehicledto.UpdateVehicleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httperr.Write(w, err)
		return
	}

	output, err := h.update.Execute(r.Context(), vehicledomain.VehicleID(id), input)
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, output)
}

func (h *VehicleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	if err := h.delete.Execute(r.Context(), vehicledomain.VehicleID(id)); err != nil {
		httperr.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
