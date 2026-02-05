package middleware

import (
	"net/http"
	"time"

	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func LoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New().String()
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			wrapped.Header().Set("X-Request-ID", requestID)

			logger.Info("входящий_запрос",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("запрос_завершён",
				"request_id", requestID,
				"status", wrapped.statusCode,
				"duration_ms", duration.Milliseconds(),
				"size", wrapped.size,
			)
		})
	}
}
