package handler

import (
	"net/http"
	"time"
)

const (
	RefreshTokenCookie = "refresh_token"
	CookiePath         = "/"
	CookieMaxAge       = 7 * 24 * 60 * 60 // 7 дней в секундах
)

// SetRefreshTokenCookie устанавливает HttpOnly cookie с refresh токеном
func SetRefreshTokenCookie(w http.ResponseWriter, token string, expiresAt time.Time, isProduction bool) {
	cookie := &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    token,
		Path:     CookiePath,
		Expires:  expiresAt,
		MaxAge:   CookieMaxAge,
		HttpOnly: true,        // Защита от XSS
		Secure:   isProduction, // HTTPS only в production
		SameSite: http.SameSiteLaxMode, // CSRF защита
	}

	http.SetCookie(w, cookie)
}

// GetRefreshTokenFromCookie извлекает refresh токен из cookie
func GetRefreshTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenCookie)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// ClearRefreshTokenCookie удаляет cookie с refresh токеном
func ClearRefreshTokenCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     CookiePath,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}
