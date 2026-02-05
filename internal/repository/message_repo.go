package repository

import (
	"context"
	"database/sql"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/google/uuid"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *domain.TaskMessage) error {
	return r.CreateWithTx(ctx, nil, message)
}

func (r *messageRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, message *domain.TaskMessage) error {
	query := `
		INSERT INTO task_messages (id, task_id, author_id, content, is_system_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, message.ID, message.TaskID, message.AuthorID,
			message.Content, message.IsSystemMessage, message.CreatedAt, message.UpdatedAt)
	} else {
		_, err = r.db.ExecContext(ctx, query, message.ID, message.TaskID, message.AuthorID,
			message.Content, message.IsSystemMessage, message.CreatedAt, message.UpdatedAt)
	}

	if err != nil {
		return errors.Internal(err, "Не удалось создать сообщение")
	}

	return nil
}

func (r *messageRepository) GetByTask(ctx context.Context, taskID uuid.UUID) ([]*domain.TaskMessage, error) {
	query := `
		SELECT id, task_id, author_id, content, is_system_message, created_at, updated_at
		FROM task_messages
		WHERE task_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить сообщения")
	}
	defer rows.Close()

	messages := []*domain.TaskMessage{}
	for rows.Next() {
		m := &domain.TaskMessage{}
		err := rows.Scan(&m.ID, &m.TaskID, &m.AuthorID, &m.Content, &m.IsSystemMessage, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, errors.Internal(err, "Не удалось обработать данные сообщения")
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.TaskMessage, error) {
	query := `
		SELECT id, task_id, author_id, content, is_system_message, created_at, updated_at
		FROM task_messages
		WHERE id = $1 AND deleted_at IS NULL
	`

	m := &domain.TaskMessage{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.TaskID, &m.AuthorID, &m.Content, &m.IsSystemMessage, &m.CreatedAt, &m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound("Сообщение не найдено")
	}
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить сообщение")
	}

	return m, nil
}

func (r *messageRepository) Update(ctx context.Context, message *domain.TaskMessage) error {
	query := `UPDATE task_messages SET content = $1 WHERE id = $2 AND deleted_at IS NULL AND is_system_message = false`

	result, err := r.db.ExecContext(ctx, query, message.Content, message.ID)
	if err != nil {
		return errors.Internal(err, "Не удалось обновить сообщение")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Сообщение не найдено или является системным")
	}

	return nil
}

func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE task_messages SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal(err, "Не удалось удалить сообщение")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Сообщение не найдено")
	}

	return nil
}
