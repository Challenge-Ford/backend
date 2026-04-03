package middleware

import (
	"net/http"

	"torque/cmd/api/httperr"
	"torque/internal/core/appctx"
	"torque/internal/core/apperr"

	"github.com/google/uuid"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(r.Header.Get("x-user-id"))
		if err != nil {
			httperr.Write(w, apperr.Unauthorized("missing or invalid x-user-id"))
			return
		}

		ctx := appctx.WithAuth(r.Context(), appctx.AuthContext{
			UserID: userID,
			Role:   r.Header.Get("x-user-role"),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
