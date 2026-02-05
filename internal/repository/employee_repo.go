package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/google/uuid"
)

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	query := `
		INSERT INTO employees (id, name, department, position, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		employee.ID,
		employee.Name,
		employee.Department,
		employee.Position,
		employee.Email,
		employee.PasswordHash,
		employee.CreatedAt,
		employee.UpdatedAt,
	)

	if err != nil {
		return errors.Internal(err, "Не удалось создать сотрудника")
	}

	return nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	query := `
		SELECT id, name, department, position, email, password_hash, created_at, updated_at, deleted_at
		FROM employees
		WHERE id = $1 AND deleted_at IS NULL
	`

	employee := &domain.Employee{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Department,
		&employee.Position,
		&employee.Email,
		&employee.PasswordHash,
		&employee.CreatedAt,
		&employee.UpdatedAt,
		&employee.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound("Сотрудник не найден")
	}
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить сотрудника")
	}

	return employee, nil
}

func (r *employeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	query := `
		SELECT id, name, department, position, email, password_hash, created_at, updated_at, deleted_at
		FROM employees
		WHERE email = $1 AND deleted_at IS NULL
	`

	employee := &domain.Employee{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Department,
		&employee.Position,
		&employee.Email,
		&employee.PasswordHash,
		&employee.CreatedAt,
		&employee.UpdatedAt,
		&employee.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound("Сотрудник не найден")
	}
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить сотрудника")
	}

	return employee, nil
}

func (r *employeeRepository) GetAll(ctx context.Context, filter EmployeeFilter) ([]*domain.Employee, int, error) {
	query := `
		SELECT id, name, department, position, email, created_at, updated_at
		FROM employees
		WHERE deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM employees WHERE deleted_at IS NULL`

	args := []interface{}{}
	argPos := 1

	if filter.Department != "" {
		query += fmt.Sprintf(" AND department = $%d", argPos)
		countQuery += fmt.Sprintf(" AND department = $%d", argPos)
		args = append(args, filter.Department)
		argPos++
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, errors.Internal(err, "Не удалось подсчитать сотрудников")
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
		return nil, 0, errors.Internal(err, "Не удалось получить список сотрудников")
	}
	defer rows.Close()

	employees := []*domain.Employee{}
	for rows.Next() {
		employee := &domain.Employee{}
		err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Department,
			&employee.Position,
			&employee.Email,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Internal(err, "Не удалось обработать данные сотрудника")
		}
		employees = append(employees, employee)
	}

	return employees, total, nil
}

func (r *employeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	query := `
		UPDATE employees
		SET name = $1, department = $2, position = $3, email = $4
		WHERE id = $5 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		employee.Name,
		employee.Department,
		employee.Position,
		employee.Email,
		employee.ID,
	)

	if err != nil {
		return errors.Internal(err, "Не удалось обновить сотрудника")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Сотрудник не найден")
	}

	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE employees
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal(err, "Не удалось удалить сотрудника")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Сотрудник не найден")
	}

	return nil
}
