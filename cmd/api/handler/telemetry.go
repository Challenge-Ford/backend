package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"torque/cmd/api/httperr"
	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
)

type telemetryLister interface {
	Execute(ctx context.Context, input telemetrydto.ListTelemetryInput) (*telemetrydto.TelemetryListOutput, error)
}

type dtcLister interface {
	Execute(ctx context.Context, vehicleID uuid.UUID) (*telemetrydto.DTCListOutput, error)
}

type TelemetryHandler struct {
	listTelemetry telemetryLister
	listDTCs      dtcLister
}

func NewTelemetryHandler(listTelemetry telemetryLister, listDTCs dtcLister) *TelemetryHandler {
	return &TelemetryHandler{listTelemetry: listTelemetry, listDTCs: listDTCs}
}

func (h *TelemetryHandler) ListTelemetry(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httperr.Write(w, apperr.BadRequest("invalid vehicle id"))
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
		after = t
	}

	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.listTelemetry.Execute(r.Context(), telemetrydto.ListTelemetryInput{
		VehicleID: id,
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
		httperr.Write(w, apperr.BadRequest("invalid vehicle id"))
		return
	}

	result, err := h.listDTCs.Execute(r.Context(), id)
	if err != nil {
		httperr.Write(w, err)
		return
	}

	httperr.JSON(w, http.StatusOK, result)
}

func parseTime(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
