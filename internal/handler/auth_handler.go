package handler

import (
	"net/http"

	"github.com/dmitry/taskmanager/internal/dto"
	"github.com/dmitry/taskmanager/internal/service"
	"github.com/dmitry/taskmanager/pkg/validator"
)

type AuthHandler struct {
	authService  *service.AuthService
	validator    *validator.Validator
	isProduction bool
}

func NewAuthHandler(authService *service.AuthService, validator *validator.Validator, isProduction bool) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		validator:    validator,
		isProduction: isProduction,
	}
}

// Register создает новую учетную запись сотрудника и возвращает токены
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	if err := h.validator.Validate(req); err != nil {
		RespondError(w, err)
		return
	}

	userAgent := r.UserAgent()
	ipAddress := getIPAddress(r)

	employee, err := h.authService.Register(r.Context(), req.Name, req.Department, req.Position, req.Email, req.Password)
	if err != nil {
		RespondError(w, err)
		return
	}

	tokens, _, err := h.authService.Login(r.Context(), req.Email, req.Password, userAgent, ipAddress)
	if err != nil {
		RespondError(w, err)
		return
	}

	// Устанавливаем refresh token в HttpOnly cookie
	SetRefreshTokenCookie(w, tokens.RefreshToken, tokens.ExpiresAt, h.isProduction)

	// Возвращаем только access token в JSON
	response := dto.ToAuthResponse(tokens.AccessToken, tokens.ExpiresAt, employee)
	RespondJSON(w, http.StatusCreated, response)
}

// Login аутентифицирует сотрудника и возвращает токены
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	if err := h.validator.Validate(req); err != nil {
		RespondError(w, err)
		return
	}

	userAgent := r.UserAgent()
	ipAddress := getIPAddress(r)

	tokens, employee, err := h.authService.Login(r.Context(), req.Email, req.Password, userAgent, ipAddress)
	if err != nil {
		RespondError(w, err)
		return
	}

	// Устанавливаем refresh token в HttpOnly cookie
	SetRefreshTokenCookie(w, tokens.RefreshToken, tokens.ExpiresAt, h.isProduction)

	// Возвращаем только access token в JSON
	response := dto.ToAuthResponse(tokens.AccessToken, tokens.ExpiresAt, employee)
	RespondJSON(w, http.StatusOK, response)
}

// RefreshToken генерирует новый токен доступа используя refresh токен из cookie
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Читаем refresh token из cookie
	refreshToken, err := GetRefreshTokenFromCookie(r)
	if err != nil {
		RespondError(w, err)
		return
	}

	userAgent := r.UserAgent()
	ipAddress := getIPAddress(r)

	tokens, err := h.authService.RefreshToken(r.Context(), refreshToken, userAgent, ipAddress)
	if err != nil {
		// При ошибке обновления удаляем cookie
		ClearRefreshTokenCookie(w)
		RespondError(w, err)
		return
	}

	// Устанавливаем новый refresh token в cookie
	SetRefreshTokenCookie(w, tokens.RefreshToken, tokens.ExpiresAt, h.isProduction)

	// Возвращаем только access token
	response := dto.ToTokenResponse(tokens.AccessToken, tokens.ExpiresAt)
	RespondJSON(w, http.StatusOK, response)
}

// Logout отзывает refresh токен из cookie
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Читаем refresh token из cookie
	refreshToken, err := GetRefreshTokenFromCookie(r)
	if err != nil {
		// Если cookie нет, просто очищаем и возвращаем успех
		ClearRefreshTokenCookie(w)
		RespondJSON(w, http.StatusOK, map[string]string{"message": "Выход выполнен успешно"})
		return
	}

	// Отзываем токен в базе данных
	if err := h.authService.Logout(r.Context(), refreshToken); err != nil {
		// Даже при ошибке удаляем cookie
		ClearRefreshTokenCookie(w)
		RespondError(w, err)
		return
	}

	// Удаляем cookie
	ClearRefreshTokenCookie(w)
	RespondJSON(w, http.StatusOK, map[string]string{"message": "Выход выполнен успешно"})
}

// getIPAddress извлекает IP-адрес из запроса
func getIPAddress(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}
