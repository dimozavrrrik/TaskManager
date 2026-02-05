package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dmitry/taskmanager/internal/service"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/google/uuid"
)

type contextKey string

const EmployeeIDKey contextKey = "employee_id"

func AuthMiddleware(jwtService *service.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, errors.Unauthorized("Требуется заголовок Authorization"))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, errors.Unauthorized("Неверный формат авторизации. Используйте: Bearer <token>"))
				return
			}

			tokenString := parts[1]

			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				respondError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), EmployeeIDKey, claims.EmployeeID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetEmployeeIDFromContext(ctx context.Context) (uuid.UUID, error) {
	employeeID, ok := ctx.Value(EmployeeIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.Unauthorized("ID сотрудника не найден в контексте")
	}
	return employeeID, nil
}

func respondError(w http.ResponseWriter, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.Internal(err, "Произошла непредвиденная ошибка")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatusCode())

	response := `{"success":false,"error":{"code":"` + string(appErr.Code) + `","message":"` + appErr.Message + `"}}`
	w.Write([]byte(response))
}
