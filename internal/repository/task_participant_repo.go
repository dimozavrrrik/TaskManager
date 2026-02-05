package repository

import (
	"context"
	"database/sql"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/google/uuid"
)

type taskParticipantRepository struct {
	db *sql.DB
}

func NewTaskParticipantRepository(db *sql.DB) TaskParticipantRepository {
	return &taskParticipantRepository{db: db}
}

func (r *taskParticipantRepository) AddParticipant(ctx context.Context, participant *domain.TaskParticipant) error {
	return r.AddParticipantWithTx(ctx, nil, participant)
}

func (r *taskParticipantRepository) AddParticipantWithTx(ctx context.Context, tx *sql.Tx, participant *domain.TaskParticipant) error {
	query := `
		INSERT INTO task_participants (id, task_id, employee_id, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (task_id, employee_id, role) DO NOTHING
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, participant.ID, participant.TaskID,
			participant.EmployeeID, participant.Role, participant.CreatedAt)
	} else {
		_, err = r.db.ExecContext(ctx, query, participant.ID, participant.TaskID,
			participant.EmployeeID, participant.Role, participant.CreatedAt)
	}

	if err != nil {
		return errors.Internal(err, "Не удалось добавить участника")
	}

	return nil
}

func (r *taskParticipantRepository) RemoveParticipant(ctx context.Context, taskID, employeeID uuid.UUID, role domain.ParticipantRole) error {
	query := `DELETE FROM task_participants WHERE task_id = $1 AND employee_id = $2 AND role = $3`

	result, err := r.db.ExecContext(ctx, query, taskID, employeeID, role)
	if err != nil {
		return errors.Internal(err, "Не удалось удалить участника")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NotFound("Участник не найден")
	}

	return nil
}

func (r *taskParticipantRepository) GetParticipants(ctx context.Context, taskID uuid.UUID) ([]*domain.TaskParticipant, error) {
	query := `
		SELECT id, task_id, employee_id, role, created_at
		FROM task_participants
		WHERE task_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить список участников")
	}
	defer rows.Close()

	participants := []*domain.TaskParticipant{}
	for rows.Next() {
		p := &domain.TaskParticipant{}
		err := rows.Scan(&p.ID, &p.TaskID, &p.EmployeeID, &p.Role, &p.CreatedAt)
		if err != nil {
			return nil, errors.Internal(err, "Не удалось обработать данные участника")
		}
		participants = append(participants, p)
	}

	return participants, nil
}

func (r *taskParticipantRepository) GetParticipantsByEmployee(ctx context.Context, employeeID uuid.UUID) ([]*domain.TaskParticipant, error) {
	query := `
		SELECT id, task_id, employee_id, role, created_at
		FROM task_participants
		WHERE employee_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, employeeID)
	if err != nil {
		return nil, errors.Internal(err, "Не удалось получить список участников")
	}
	defer rows.Close()

	participants := []*domain.TaskParticipant{}
	for rows.Next() {
		p := &domain.TaskParticipant{}
		err := rows.Scan(&p.ID, &p.TaskID, &p.EmployeeID, &p.Role, &p.CreatedAt)
		if err != nil {
			return nil, errors.Internal(err, "Не удалось обработать данные участника")
		}
		participants = append(participants, p)
	}

	return participants, nil
}
