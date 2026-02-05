package repository

import (
	"context"
	"database/sql"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/google/uuid"
)

type EmployeeFilter struct {
	Department string
	Page       int
	PageSize   int
}

type EmployeeRepository interface {
	Create(ctx context.Context, employee *domain.Employee) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	GetByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetAll(ctx context.Context, filter EmployeeFilter) ([]*domain.Employee, int, error)
	Update(ctx context.Context, employee *domain.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TaskFilter struct {
	Status     []domain.TaskStatus
	Priority   *int
	Archived   *bool
	EmployeeID *uuid.UUID
	Page       int
	PageSize   int
}

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	CreateWithTx(ctx context.Context, tx *sql.Tx, task *domain.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetAll(ctx context.Context, filter TaskFilter) ([]*domain.Task, int, error)
	Update(ctx context.Context, task *domain.Task) error
	UpdateWithTx(ctx context.Context, tx *sql.Tx, task *domain.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) (domain.TaskStatus, error)
	Archive(ctx context.Context, id uuid.UUID) error
	GetTasksForEmployee(ctx context.Context, employeeID uuid.UUID, filter TaskFilter) ([]*domain.Task, int, error)
}

type TaskParticipantRepository interface {
	AddParticipant(ctx context.Context, participant *domain.TaskParticipant) error
	AddParticipantWithTx(ctx context.Context, tx *sql.Tx, participant *domain.TaskParticipant) error
	RemoveParticipant(ctx context.Context, taskID, employeeID uuid.UUID, role domain.ParticipantRole) error
	GetParticipants(ctx context.Context, taskID uuid.UUID) ([]*domain.TaskParticipant, error)
	GetParticipantsByEmployee(ctx context.Context, employeeID uuid.UUID) ([]*domain.TaskParticipant, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *domain.TaskMessage) error
	CreateWithTx(ctx context.Context, tx *sql.Tx, message *domain.TaskMessage) error
	GetByTask(ctx context.Context, taskID uuid.UUID) ([]*domain.TaskMessage, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TaskMessage, error)
	Update(ctx context.Context, message *domain.TaskMessage) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TimeEntryFilter struct {
	StartDate *string
	EndDate   *string
	Page      int
	PageSize  int
}

type TimeEntryRepository interface {
	Create(ctx context.Context, entry *domain.TimeEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TimeEntry, error)
	GetByTask(ctx context.Context, taskID uuid.UUID) ([]*domain.TimeEntry, error)
	GetByEmployee(ctx context.Context, employeeID uuid.UUID, filter TimeEntryFilter) ([]*domain.TimeEntry, error)
	Update(ctx context.Context, entry *domain.TimeEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTaskTimeSummary(ctx context.Context, taskID uuid.UUID) (*domain.TimeSummary, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	RevokeByTokenHash(ctx context.Context, tokenHash string) error
	RevokeAllByEmployee(ctx context.Context, employeeID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
