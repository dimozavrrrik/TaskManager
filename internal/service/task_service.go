package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/internal/repository"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo        repository.TaskRepository
	participantRepo repository.TaskParticipantRepository
	messageRepo     repository.MessageRepository
	employeeRepo    repository.EmployeeRepository
	db              *sql.DB
	logger          *logger.Logger
}

func NewTaskService(
	taskRepo repository.TaskRepository,
	participantRepo repository.TaskParticipantRepository,
	messageRepo repository.MessageRepository,
	employeeRepo repository.EmployeeRepository,
	db *sql.DB,
	logger *logger.Logger,
) *TaskService {
	return &TaskService{
		taskRepo:        taskRepo,
		participantRepo: participantRepo,
		messageRepo:     messageRepo,
		employeeRepo:    employeeRepo,
		db:              db,
		logger:          logger,
	}
}

type CreateTaskRequest struct {
	Title        string
	Description  string
	Priority     int
	CreatedBy    uuid.UUID
	DueDate      *string
	Participants []ParticipantInput
}

type ParticipantInput struct {
	EmployeeID uuid.UUID
	Role       domain.ParticipantRole
}

func (s *TaskService) CreateTask(ctx context.Context, req CreateTaskRequest) (*domain.Task, error) {
	if _, err := s.employeeRepo.GetByID(ctx, req.CreatedBy); err != nil {
		return nil, errors.BadRequest("Сотрудник-создатель не найден")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Internal(err, "Не удалось начать транзакцию")
	}
	defer tx.Rollback()

	task := domain.NewTask(req.Title, req.Description, req.Priority, req.CreatedBy, nil)

	if err := s.taskRepo.CreateWithTx(ctx, tx, task); err != nil {
		return nil, err
	}

	for _, p := range req.Participants {
		if _, err := s.employeeRepo.GetByID(ctx, p.EmployeeID); err != nil {
			return nil, errors.BadRequest(fmt.Sprintf("Сотрудник %s не найден", p.EmployeeID))
		}

		participant := domain.NewTaskParticipant(task.ID, p.EmployeeID, p.Role)
		if err := s.participantRepo.AddParticipantWithTx(ctx, tx, participant); err != nil {
			return nil, err
		}
	}

	systemMsg := domain.NewSystemMessage(task.ID, "Задача создана")
	if err := s.messageRepo.CreateWithTx(ctx, tx, systemMsg); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Internal(err, "Не удалось зафиксировать транзакцию")
	}

	s.logger.Info("Задача создана", "task_id", task.ID, "created_by", req.CreatedBy)

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *TaskService) GetAllTasks(ctx context.Context, filter repository.TaskFilter) ([]*domain.Task, int, error) {
	return s.taskRepo.GetAll(ctx, filter)
}

func (s *TaskService) GetTasksForEmployee(ctx context.Context, employeeID uuid.UUID, filter repository.TaskFilter) ([]*domain.Task, int, error) {
	return s.taskRepo.GetTasksForEmployee(ctx, employeeID, filter)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domain.Task) error {
	return s.taskRepo.Update(ctx, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, newStatus domain.TaskStatus) error {
	if !newStatus.IsValid() {
		return errors.BadRequest("Неверный статус задачи")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Internal(err, "Не удалось начать транзакцию")
	}
	defer tx.Rollback()

	oldStatus, err := s.taskRepo.UpdateStatus(ctx, taskID, newStatus)
	if err != nil {
		return err
	}

	if oldStatus != newStatus {
		content := fmt.Sprintf("Статус задачи изменён с '%s' на '%s'", oldStatus, newStatus)
		systemMsg := domain.NewSystemMessage(taskID, content)
		if err := s.messageRepo.CreateWithTx(ctx, tx, systemMsg); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Internal(err, "Не удалось зафиксировать транзакцию")
	}

	s.logger.Info("Статус задачи обновлён", "task_id", taskID, "old_status", oldStatus, "new_status", newStatus)

	return nil
}

func (s *TaskService) ArchiveTask(ctx context.Context, id uuid.UUID) error {
	if err := s.taskRepo.Archive(ctx, id); err != nil {
		return err
	}

	s.logger.Info("Задача архивирована", "task_id", id)

	return nil
}

func (s *TaskService) AddParticipant(ctx context.Context, taskID, employeeID uuid.UUID, role domain.ParticipantRole) error {
	if !role.IsValid() {
		return errors.BadRequest("Неверная роль участника")
	}

	if _, err := s.taskRepo.GetByID(ctx, taskID); err != nil {
		return err
	}

	if _, err := s.employeeRepo.GetByID(ctx, employeeID); err != nil {
		return errors.BadRequest("Сотрудник не найден")
	}

	participant := domain.NewTaskParticipant(taskID, employeeID, role)
	return s.participantRepo.AddParticipant(ctx, participant)
}

func (s *TaskService) RemoveParticipant(ctx context.Context, taskID, employeeID uuid.UUID, role domain.ParticipantRole) error {
	return s.participantRepo.RemoveParticipant(ctx, taskID, employeeID, role)
}

func (s *TaskService) GetParticipants(ctx context.Context, taskID uuid.UUID) ([]*domain.TaskParticipant, error) {
	return s.participantRepo.GetParticipants(ctx, taskID)
}
