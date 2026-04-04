package middleware

import (
	"net/http"

	"go.uber.org/zap"
	coremiddleware "torque/internal/core/middleware"
)

func Logger(log *zap.Logger) func(http.Handler) http.Handler {
	return coremiddleware.Logger(log)
}
