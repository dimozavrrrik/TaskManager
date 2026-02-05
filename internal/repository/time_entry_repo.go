package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/google/uuid"
)

type timeEntryRepository struct {
	db *sql.DB
}

func NewTimeEntryRepository(db *sql.DB) TimeEntryRepository {
	return &timeEntryRepository{db: db}
}

func (r *timeEntryRepository) Create(ctx context.Context, entry *domain.TimeEntry) error {
	query := `
		INSERT INTO time_entries (id, task_id, employee_id, hours, description, entry_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query, entry.ID, entry.TaskID, entry.EmployeeID,
		entry.Hours, entry.Description, entry.EntryDate, entry.CreatedAt, entry.UpdatedAt)

	if err != nil {
		return errors.Internal(err, "Не удалось создать запись времени")
	}

	return nil
}

func (r *timeEntryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.TimeEntry, error) {
	query := `
		SELECT id, task_id, employee_id, hours, description, entry_date, created_at, updated_at
		FROM time_entries
		WHERE id = $1 AND deleted_at IS NULL
	`

	entry := &domain.TimeEntry{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID, &entry.TaskID, &entry.EmployeeID, &entry.Hours,
		&entry.Description, &entry.EntryDate, &entry.CreatedAt, &entry.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound("Запись времени не найдена")
	}
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить запись времени")
	}

	return entry, nil
}

func (r *timeEntryRepository) GetByTask(ctx context.Context, taskID uuid.UUID) ([]*domain.TimeEntry, error) {
	query := `
		SELECT id, task_id, employee_id, hours, description, entry_date, created_at, updated_at
		FROM time_entries
		WHERE task_id = $1 AND deleted_at IS NULL
		ORDER BY entry_date DESC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить записи времени")
	}
	defer rows.Close()

	entries := []*domain.TimeEntry{}
	for rows.Next() {
		entry := &domain.TimeEntry{}
		err := rows.Scan(&entry.ID, &entry.TaskID, &entry.EmployeeID, &entry.Hours,
			&entry.Description, &entry.EntryDate, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, errors.Internal(err, "Не удалось обработать запись времени")
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *timeEntryRepository) GetByEmployee(ctx context.Context, employeeID uuid.UUID, filter TimeEntryFilter) ([]*domain.TimeEntry, error) {
	query := `
		SELECT id, task_id, employee_id, hours, description, entry_date, created_at, updated_at
		FROM time_entries
		WHERE employee_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{employeeID}
	argPos := 2

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND entry_date >= $%d", argPos)
		args = append(args, *filter.StartDate)
		argPos++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND entry_date <= $%d", argPos)
		args = append(args, *filter.EndDate)
		argPos++
	}

	query += " ORDER BY entry_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить записи времени")
	}
	defer rows.Close()

	entries := []*domain.TimeEntry{}
	for rows.Next() {
		entry := &domain.TimeEntry{}
		err := rows.Scan(&entry.ID, &entry.TaskID, &entry.EmployeeID, &entry.Hours,
			&entry.Description, &entry.EntryDate, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, errors.Internal(err, "Не удалось обработать запись времени")
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *timeEntryRepository) Update(ctx context.Context, entry *domain.TimeEntry) error {
	query := `
		UPDATE time_entries
		SET hours = $1, description = $2, entry_date = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, entry.Hours, entry.Description, entry.EntryDate, entry.ID)
	if err != nil {
		return errors.Internal(err, "Не удалось обновить запись времени")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Запись времени не найдена")
	}

	return nil
}

func (r *timeEntryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE time_entries SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal(err, "Не удалось удалить запись времени")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Запись времени не найдена")
	}

	return nil
}

func (r *timeEntryRepository) GetTaskTimeSummary(ctx context.Context, taskID uuid.UUID) (*domain.TimeSummary, error) {
	query := `
		SELECT task_id, total_hours, entry_count, unique_employees
		FROM task_time_summary
		WHERE task_id = $1
	`

	summary := &domain.TimeSummary{}
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&summary.TaskID, &summary.TotalHours, &summary.EntryCount, &summary.UniqueEmployees,
	)

	if err == sql.ErrNoRows {
		return &domain.TimeSummary{
			TaskID:          taskID,
			TotalHours:      0,
			EntryCount:      0,
			UniqueEmployees: 0,
		}, nil
	}
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить сводку по времени")
	}

	return summary, nil
}
