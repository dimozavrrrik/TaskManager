package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dmitry/taskmanager/pkg/logger"
)

func RecoveryMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("паника_перехвачена",
						"error", fmt.Sprintf("%v", err),
						"stack", string(debug.Stack()),
						"path", r.URL.Path,
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"success":false,"error":{"code":"INTERNAL_ERROR","message":"Внутренняя ошибка сервера"}}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
