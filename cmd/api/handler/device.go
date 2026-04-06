package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"torque/cmd/api/httperr"
	"torque/internal/core/apperr"
	"torque/internal/core/pagination"
	devicedto "torque/internal/modules/device/application/dto"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicedomain "torque/internal/modules/device/domain"
)

type DeviceHandler struct {
	list         *deviceusecase.ListDevicesUseCase
	create       *deviceusecase.CreateDeviceUseCase
	commission   *deviceusecase.CommissionDeviceUseCase
	decommission *deviceusecase.DecommissionDeviceUseCase
}

func NewDeviceHandler(
	list *deviceusecase.ListDevicesUseCase,
	create *deviceusecase.CreateDeviceUseCase,
	commission *deviceusecase.CommissionDeviceUseCase,
	decommission *deviceusecase.DecommissionDeviceUseCase,
) *DeviceHandler {
	return &DeviceHandler{list: list, create: create, commission: commission, decommission: decommission}
}

func (h *DeviceHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))

	result, err := h.list.Execute(r.Context(), pagination.Page{Page: page, PerPage: perPage})
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, result)
}

func (h *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input devicedto.CreateDeviceInput
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

func (h *DeviceHandler) Commission(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, apperr.BadRequest("invalid device id"))
		return
	}

	var input devicedto.CommissionDeviceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httperr.Write(w, err)
		return
	}

	output, err := h.commission.Execute(r.Context(), devicedomain.DeviceID(id), input)
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, output)
}

func (h *DeviceHandler) Decommission(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, apperr.BadRequest("invalid device id"))
		return
	}

	output, err := h.decommission.Execute(r.Context(), devicedomain.DeviceID(id))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, output)
}

