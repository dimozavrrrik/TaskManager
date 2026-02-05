package service

import (
	"context"
	"time"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/internal/repository"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	employeeRepo     repository.EmployeeRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtService       *JWTService
	logger           *logger.Logger
}

func NewAuthService(
	employeeRepo repository.EmployeeRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtService *JWTService,
	logger *logger.Logger,
) *AuthService {
	return &AuthService{
		employeeRepo:     employeeRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		logger:           logger,
	}
}

type AuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Register создает нового сотрудника с паролем
func (s *AuthService) Register(ctx context.Context, name, department, position, email, password string) (*domain.Employee, error) {
	existing, err := s.employeeRepo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, errors.Conflict("Email уже зарегистрирован")
	}

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	employee := domain.NewEmployee(name, department, position, email)
	employee.PasswordHash = hashedPassword

	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		return nil, err
	}

	s.logger.Info("Сотрудник зарегистрирован", "employee_id", employee.ID, "email", email)

	return employee, nil
}

// Login аутентифицирует пользователя и возвращает токены
func (s *AuthService) Login(ctx context.Context, email, password, userAgent, ipAddress string) (*AuthTokens, *domain.Employee, error) {
	employee, err := s.employeeRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, errors.Unauthorized("Неверный email или пароль")
	}

	if employee.PasswordHash == "" {
		return nil, nil, errors.Unauthorized("Пароль не установлен для этой учётной записи")
	}

	if err := s.verifyPassword(employee.PasswordHash, password); err != nil {
		return nil, nil, errors.Unauthorized("Неверный email или пароль")
	}

	tokens, err := s.generateTokens(ctx, employee, userAgent, ipAddress)
	if err != nil {
		return nil, nil, err
	}

	s.logger.Info("Сотрудник вошёл в систему", "employee_id", employee.ID, "email", email)

	return tokens, employee, nil
}

// RefreshToken генерирует новый токен доступа из refresh токена
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, userAgent, ipAddress string) (*AuthTokens, error) {
	employeeID, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	tokenHash := s.jwtService.HashToken(refreshToken)
	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	if !storedToken.IsValid() {
		return nil, errors.Unauthorized("Refresh-токен недействителен или истёк")
	}

	employee, err := s.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, errors.Unauthorized("Сотрудник не найден")
	}

	if err := s.refreshTokenRepo.RevokeByTokenHash(ctx, tokenHash); err != nil {
		s.logger.Error("Не удалось отозвать старый refresh-токен", "error", err)
	}

	tokens, err := s.generateTokens(ctx, employee, userAgent, ipAddress)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Токены обновлены", "employee_id", employee.ID)

	return tokens, nil
}

// Logout отзывает refresh токен
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := s.jwtService.HashToken(refreshToken)

	if err := s.refreshTokenRepo.RevokeByTokenHash(ctx, tokenHash); err != nil {
		s.logger.Warn("Не удалось отозвать refresh-токен при выходе", "error", err)
		return nil
	}

	s.logger.Info("Сотрудник вышел из системы")

	return nil
}

// LogoutAll отзывает все refresh токены для сотрудника
func (s *AuthService) LogoutAll(ctx context.Context, employeeID uuid.UUID) error {
	if err := s.refreshTokenRepo.RevokeAllByEmployee(ctx, employeeID); err != nil {
		return err
	}

	s.logger.Info("Все сессии отозваны", "employee_id", employeeID)

	return nil
}

// generateTokens создает токен доступа и refresh токен
func (s *AuthService) generateTokens(ctx context.Context, employee *domain.Employee, userAgent, ipAddress string) (*AuthTokens, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(employee.ID, employee.Email, employee.Name)
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := s.jwtService.GenerateRefreshToken(employee.ID)
	if err != nil {
		return nil, err
	}

	tokenHash := s.jwtService.HashToken(refreshToken)
	dbToken := domain.NewRefreshToken(employee.ID, tokenHash, expiresAt, userAgent, ipAddress)

	if err := s.refreshTokenRepo.Create(ctx, dbToken); err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// hashPassword хеширует пароль используя bcrypt
func (s *AuthService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Internal(err, "Не удалось хешировать пароль")
	}
	return string(hashedBytes), nil
}

// verifyPassword проверяет пароль на соответствие его хешу
func (s *AuthService) verifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.Unauthorized("Неверный пароль")
	}
	return nil
}
