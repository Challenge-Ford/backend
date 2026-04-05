package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"torque/cmd/api/httperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type TelemetryHandler struct {
	listTelemetry *telemetryusecase.ListTelemetryUseCase
	listDTCs      *telemetryusecase.ListActiveDTCsUseCase
}

func NewTelemetryHandler(
	listTelemetry *telemetryusecase.ListTelemetryUseCase,
	listDTCs *telemetryusecase.ListActiveDTCsUseCase,
) *TelemetryHandler {
	return &TelemetryHandler{listTelemetry: listTelemetry, listDTCs: listDTCs}
}

func (h *TelemetryHandler) ListTelemetry(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	q := r.URL.Query()

	from, err := parseTime(q.Get("from"))
	if err != nil {
		httperr.Write(w, err)
		return
	}
	to, err := parseTime(q.Get("to"))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	var after *time.Time
	if raw := q.Get("after"); raw != "" {
		t, err := parseTime(raw)
		if err != nil {
			httperr.Write(w, err)
			return
		}
		after = &t
	}

	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.listTelemetry.Execute(r.Context(), telemetrydto.ListTelemetryInput{
		VehicleID: vehicledomain.VehicleID(id),
		From:      from,
		To:        to,
		Limit:     limit,
		After:     after,
	})
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, result)
}

func (h *TelemetryHandler) ListDTCs(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	result, err := h.listDTCs.Execute(r.Context(), vehicledomain.VehicleID(id))
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, result)
}

func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}
