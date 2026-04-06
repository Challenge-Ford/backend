package handler

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"torque/cmd/api/httperr"
)

type HealthHandler struct {
	db   *pgxpool.Pool
	tsDB *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool, tsDB *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db, tsDB: tsDB}
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	httperr.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5000000000)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		httperr.JSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "db_unhealthy",
			"error":  err.Error(),
		})
		return
	}

	if err := h.tsDB.Ping(ctx); err != nil {
		httperr.JSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "tsdb_unhealthy",
			"error":  err.Error(),
		})
		return
	}

	httperr.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
