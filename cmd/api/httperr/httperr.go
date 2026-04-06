package httperr

import (
	"encoding/json"
	"errors"
	"net/http"

	"torque/internal/core/apperr"
)

type response struct {
	Error   string                `json:"error"`
	Message string                `json:"message"`
	Fields  []apperr.ValidationItem `json:"fields,omitempty"`
}

func Write(w http.ResponseWriter, err error) {
	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		r := response{
			Error:   appErr.Kind.String(),
			Message: appErr.Message,
		}
		if len(appErr.ValidationErrors) > 0 {
			r.Fields = appErr.ValidationErrors
		}
		writeJSON(w, appErr.Kind.HTTPStatus(), r)
		return
	}

	writeJSON(w, http.StatusInternalServerError, response{
		Error:   "INTERNAL_ERROR",
		Message: "internal server error",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func JSON(w http.ResponseWriter, status int, v any) {
	writeJSON(w, status, v)
}
