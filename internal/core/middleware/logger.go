package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.New().String()

			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			w.Header().Set("x-request-id", requestID)

			next.ServeHTTP(rw, r)

			fields := []zap.Field{
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.status),
				zap.Int("bytes", rw.size),
				zap.Duration("latency", time.Since(start)),
			}

			if r.URL.RawQuery != "" {
				fields = append(fields, zap.String("query", r.URL.RawQuery))
			}
			if userID := r.Header.Get("x-user-id"); userID != "" {
				fields = append(fields, zap.String("user_id", userID))
			}

			if rw.status >= 500 {
				log.Error("request", fields...)
			} else if rw.status >= 400 {
				log.Warn("request", fields...)
			} else {
				log.Info("request", fields...)
			}
		})
	}
}
