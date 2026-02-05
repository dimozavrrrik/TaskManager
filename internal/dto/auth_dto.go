package dto

import (
	"time"

	"github.com/dmitry/taskmanager/internal/domain"
)

// RegisterRequest - запрос на регистрацию пользователя
type RegisterRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Department string `json:"department" validate:"required,max=100"`
	Position   string `json:"position" validate:"required,max=100"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8,max=72"`
}

// LoginRequest - запрос на вход пользователя
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest больше не нужен - токен читается из cookie

// LogoutRequest больше не нужен - токен читается из cookie

// AuthResponse - ответ с токеном доступа и информацией о сотруднике (refresh токен в HttpOnly cookie)
type AuthResponse struct {
	AccessToken string           `json:"access_token"`
	ExpiresAt   time.Time        `json:"expires_at"`
	Employee    EmployeeResponse `json:"employee"`
}

// TokenResponse - ответ только с токеном доступа (refresh токен в HttpOnly cookie)
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// ToAuthResponse преобразует данные сервиса в DTO AuthResponse
func ToAuthResponse(accessToken string, expiresAt time.Time, employee *domain.Employee) AuthResponse {
	return AuthResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		Employee:    ToEmployeeResponse(employee),
	}
}

// ToTokenResponse преобразует данные сервиса в DTO TokenResponse
func ToTokenResponse(accessToken string, expiresAt time.Time) TokenResponse {
	return TokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}
}
