package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/google/uuid"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	return r.CreateWithTx(ctx, nil, task)
}

func (r *taskRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, priority, created_by, archived, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, task.ID, task.Title, task.Description, task.Status,
			task.Priority, task.CreatedBy, task.Archived, task.DueDate, task.CreatedAt, task.UpdatedAt)
	} else {
		_, err = r.db.ExecContext(ctx, query, task.ID, task.Title, task.Description, task.Status,
			task.Priority, task.CreatedBy, task.Archived, task.DueDate, task.CreatedAt, task.UpdatedAt)
	}

	if err != nil {
		return errors.Internal(err, "Не удалось создать задачу")
	}

	return nil
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, created_by, archived, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL
	`

	task := &domain.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&task.CreatedBy, &task.Archived, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound("Задача не найдена")
	}
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить задачу")
	}

	return task, nil
}

func (r *taskRepository) GetAll(ctx context.Context, filter TaskFilter) ([]*domain.Task, int, error) {
	query := `SELECT id, title, description, status, priority, created_by, archived, due_date, created_at, updated_at FROM tasks WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL`

	args := []interface{}{}
	argPos := 1

	if len(filter.Status) > 0 {
		placeholders := []string{}
		for _, status := range filter.Status {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argPos))
			args = append(args, status)
			argPos++
		}
		statusFilter := " AND status IN (" + strings.Join(placeholders, ",") + ")"
		query += statusFilter
		countQuery += statusFilter
	}

	if filter.Priority != nil {
		query += fmt.Sprintf(" AND priority = $%d", argPos)
		countQuery += fmt.Sprintf(" AND priority = $%d", argPos)
		args = append(args, *filter.Priority)
		argPos++
	}

	if filter.Archived != nil {
		query += fmt.Sprintf(" AND archived = $%d", argPos)
		countQuery += fmt.Sprintf(" AND archived = $%d", argPos)
		args = append(args, *filter.Archived)
		argPos++
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, errors.Internal(err, "Не удалось подсчитать задачи")
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.Internal(err, "Не удалось получить список задач")
	}
	defer rows.Close()

	tasks := []*domain.Task{}
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
			&task.CreatedBy, &task.Archived, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, 0, errors.Internal(err, "Не удалось обработать данные задачи")
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	return r.UpdateWithTx(ctx, nil, task)
}

func (r *taskRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, priority = $4, due_date = $5
		WHERE id = $6 AND deleted_at IS NULL
	`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(ctx, query, task.Title, task.Description, task.Status,
			task.Priority, task.DueDate, task.ID)
	} else {
		result, err = r.db.ExecContext(ctx, query, task.Title, task.Description, task.Status,
			task.Priority, task.DueDate, task.ID)
	}

	if err != nil {
		return errors.Internal(err, "Не удалось обновить задачу")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Задача не найдена")
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal(err, "Не удалось удалить задачу")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Задача не найдена")
	}

	return nil
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) (domain.TaskStatus, error) {
	var oldStatus domain.TaskStatus

	querySelect := `SELECT status FROM tasks WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, querySelect, id).Scan(&oldStatus)
	if err == sql.ErrNoRows {
		return "", errors.NotFound("Задача не найдена")
	}
	if err != nil {
		return "", errors.Internal(err, "Не удалось получить статус задачи")
	}

	queryUpdate := `UPDATE tasks SET status = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err = r.db.ExecContext(ctx, queryUpdate, status, id)
	if err != nil {
		return "", errors.Internal(err, "Не удалось обновить статус задачи")
	}

	return oldStatus, nil
}

func (r *taskRepository) Archive(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tasks SET archived = true WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal(err, "Не удалось архивировать задачу")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Задача не найдена")
	}

	return nil
}

func (r *taskRepository) GetTasksForEmployee(ctx context.Context, employeeID uuid.UUID, filter TaskFilter) ([]*domain.Task, int, error) {
	query := `
		SELECT DISTINCT t.id, t.title, t.description, t.status, t.priority, t.created_by, t.archived, t.due_date, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN task_participants tp ON t.id = tp.task_id
		WHERE t.deleted_at IS NULL AND tp.employee_id = $1
	`
	countQuery := `
		SELECT COUNT(DISTINCT t.id)
		FROM tasks t
		INNER JOIN task_participants tp ON t.id = tp.task_id
		WHERE t.deleted_at IS NULL AND tp.employee_id = $1
	`

	args := []interface{}{employeeID}
	argPos := 2

	if len(filter.Status) > 0 {
		placeholders := []string{}
		for _, status := range filter.Status {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argPos))
			args = append(args, status)
			argPos++
		}
		statusFilter := " AND t.status IN (" + strings.Join(placeholders, ",") + ")"
		query += statusFilter
		countQuery += statusFilter
	}

	if filter.Archived != nil {
		query += fmt.Sprintf(" AND t.archived = $%d", argPos)
		countQuery += fmt.Sprintf(" AND t.archived = $%d", argPos)
		args = append(args, *filter.Archived)
		argPos++
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, errors.Internal(err, "Не удалось подсчитать задачи")
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" ORDER BY t.created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.Internal(err, "Не удалось получить список задач")
	}
	defer rows.Close()

	tasks := []*domain.Task{}
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
			&task.CreatedBy, &task.Archived, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, 0, errors.Internal(err, "Не удалось обработать данные задачи")
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}
