package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"torque/cmd/api/httperr"
	"torque/internal/core/appctx"
	"torque/internal/core/apperr"

	"github.com/google/uuid"
)

type AuthConfig struct {
	BypassEnabled bool
	BypassAuth    appctx.AuthContext
}

func NewAuthConfigFromEnv() (AuthConfig, error) {
	cfg := AuthConfig{}

	enabledRaw := strings.TrimSpace(os.Getenv("AUTH_BYPASS_ENABLED"))
	if enabledRaw == "" {
		return cfg, nil
	}

	enabled, err := strconv.ParseBool(enabledRaw)
	if err != nil {
		return AuthConfig{}, apperr.BadRequest("AUTH_BYPASS_ENABLED must be a boolean")
	}
	if !enabled {
		return cfg, nil
	}

	if strings.TrimSpace(os.Getenv("APP_ENV")) != "development" {
		return AuthConfig{}, apperr.BadRequest("AUTH_BYPASS_ENABLED requires APP_ENV=development")
	}

	userID, err := uuid.Parse(strings.TrimSpace(os.Getenv("AUTH_BYPASS_USER_ID")))
	if err != nil {
		return AuthConfig{}, apperr.BadRequest("AUTH_BYPASS_USER_ID must be a valid UUID")
	}

	role := strings.TrimSpace(os.Getenv("AUTH_BYPASS_USER_ROLE"))
	if role == "" {
		return AuthConfig{}, apperr.BadRequest("AUTH_BYPASS_USER_ROLE is required when AUTH_BYPASS_ENABLED=true")
	}

	cfg.BypassEnabled = true
	cfg.BypassAuth = appctx.AuthContext{
		UserID: userID,
		Role:   role,
	}

	return cfg, nil
}

func Auth(cfg AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth, ok := authFromHeaders(r)
			if !ok {
				if cfg.BypassEnabled {
					auth = cfg.BypassAuth
				} else {
					httperr.Write(w, apperr.Unauthorized("missing or invalid x-user-id"))
					return
				}
			}

			ctx := appctx.WithAuth(r.Context(), auth)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func authFromHeaders(r *http.Request) (appctx.AuthContext, bool) {
	userID, err := uuid.Parse(r.Header.Get("x-user-id"))
	if err != nil {
		return appctx.AuthContext{}, false
	}

	return appctx.AuthContext{
		UserID: userID,
		Role:   r.Header.Get("x-user-role"),
	}, true
}
