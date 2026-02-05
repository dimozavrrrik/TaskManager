package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret            []byte
	accessExpiryMin   int
	refreshExpiryDays int
}

func NewJWTService(secret string, accessExpiryMin, refreshExpiryDays int) *JWTService {
	return &JWTService{
		secret:            []byte(secret),
		accessExpiryMin:   accessExpiryMin,
		refreshExpiryDays: refreshExpiryDays,
	}
}

// GenerateAccessToken создает краткосрочный токен доступа
func (s *JWTService) GenerateAccessToken(employeeID uuid.UUID, email, name string) (string, error) {
	expiresAt := time.Now().Add(time.Duration(s.accessExpiryMin) * time.Minute)

	claims := JWTClaims{
		EmployeeID: employeeID,
		Email:      email,
		Name:       name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "taskmanager",
			Subject:   employeeID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", errors.Internal(err, "Не удалось сгенерировать токен доступа")
	}

	return signedToken, nil
}

// GenerateRefreshToken создает долгосрочный refresh токен
func (s *JWTService) GenerateRefreshToken(employeeID uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(s.refreshExpiryDays) * 24 * time.Hour)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "taskmanager",
		Subject:   employeeID.String(),
		ID:        uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, errors.Internal(err, "Не удалось сгенерировать refresh-токен")
	}

	return signedToken, expiresAt, nil
}

// ValidateAccessToken проверяет и разбирает токен доступа
func (s *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Unauthorized("Неверный метод подписи токена")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, errors.Unauthorized("Недействительный или просроченный токен")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.Unauthorized("Неверные данные токена")
}

// ValidateRefreshToken проверяет refresh токен и возвращает ID сотрудника
func (s *JWTService) ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Unauthorized("Неверный метод подписи токена")
		}
		return s.secret, nil
	})

	if err != nil {
		return uuid.Nil, errors.Unauthorized("Недействительный или просроченный refresh-токен")
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		employeeID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, errors.Unauthorized("Неверный субъект токена")
		}
		return employeeID, nil
	}

	return uuid.Nil, errors.Unauthorized("Неверные данные refresh-токена")
}

// HashToken создает SHA-256 хеш токена для безопасного хранения
func (s *JWTService) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
