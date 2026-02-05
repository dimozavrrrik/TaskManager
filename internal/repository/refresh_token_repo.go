package repository

import (
	"context"
	"database/sql"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/google/uuid"
)

type refreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, employee_id, token_hash, expires_at, created_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.EmployeeID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
		token.UserAgent,
		token.IPAddress,
	)

	if err != nil {
		return errors.Internal(err, "Не удалось создать refresh-токен")
	}

	return nil
}

func (r *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, employee_id, token_hash, expires_at, created_at, revoked_at, user_agent, ip_address
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	token := &domain.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.EmployeeID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
		&token.UserAgent,
		&token.IPAddress,
	)

	if err == sql.ErrNoRows {
		return nil, errors.Unauthorized("Недействительный refresh-токен")
	}
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить refresh-токен")
	}

	return token, nil
}

func (r *refreshTokenRepository) RevokeByTokenHash(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE token_hash = $1 AND revoked_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, tokenHash)
	if err != nil {
		return errors.Internal(err, "Не удалось отозвать refresh-токен")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Refresh-токен не найден или уже отозван")
	}

	return nil
}

func (r *refreshTokenRepository) RevokeAllByEmployee(ctx context.Context, employeeID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE employee_id = $1 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, employeeID)
	if err != nil {
		return errors.Internal(err, "Не удалось отозвать токены сотрудника")
	}

	return nil
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < CURRENT_TIMESTAMP OR revoked_at < CURRENT_TIMESTAMP - INTERVAL '30 days'
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return errors.Internal(err, "Не удалось удалить просроченные токены")
	}

	return nil
}
