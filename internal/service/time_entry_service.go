package service

import (
	"context"
	"time"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/internal/repository"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/google/uuid"
)

type TimeEntryService struct {
	repo     repository.TimeEntryRepository
	taskRepo repository.TaskRepository
	logger   *logger.Logger
}

func NewTimeEntryService(repo repository.TimeEntryRepository, taskRepo repository.TaskRepository, logger *logger.Logger) *TimeEntryService {
	return &TimeEntryService{
		repo:     repo,
		taskRepo: taskRepo,
		logger:   logger,
	}
}

func (s *TimeEntryService) CreateTimeEntry(ctx context.Context, taskID, employeeID uuid.UUID, hours float64, description, entryDateStr string) (*domain.TimeEntry, error) {
	if hours <= 0 {
		return nil, errors.BadRequest("Количество часов должно быть больше 0")
	}

	if _, err := s.taskRepo.GetByID(ctx, taskID); err != nil {
		return nil, errors.BadRequest("Задача не найдена")
	}

	entryDate, err := time.Parse("2006-01-02", entryDateStr)
	if err != nil {
		return nil, errors.BadRequest("Неверный формат даты, ожидается ГГГГ-ММ-ДД")
	}

	entry := domain.NewTimeEntry(taskID, employeeID, hours, description, entryDate)

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}

	s.logger.Info("Запись времени создана", "entry_id", entry.ID, "task_id", taskID, "hours", hours)

	return entry, nil
}

func (s *TimeEntryService) GetTimeEntry(ctx context.Context, id uuid.UUID) (*domain.TimeEntry, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TimeEntryService) GetTaskTimeEntries(ctx context.Context, taskID uuid.UUID) ([]*domain.TimeEntry, error) {
	return s.repo.GetByTask(ctx, taskID)
}

func (s *TimeEntryService) GetEmployeeTimeEntries(ctx context.Context, employeeID uuid.UUID, filter repository.TimeEntryFilter) ([]*domain.TimeEntry, error) {
	return s.repo.GetByEmployee(ctx, employeeID, filter)
}

func (s *TimeEntryService) GetTaskTimeSummary(ctx context.Context, taskID uuid.UUID) (*domain.TimeSummary, error) {
	return s.repo.GetTaskTimeSummary(ctx, taskID)
}

func (s *TimeEntryService) UpdateTimeEntry(ctx context.Context, entry *domain.TimeEntry) error {
	if entry.Hours <= 0 {
		return errors.BadRequest("Количество часов должно быть больше 0")
	}

	return s.repo.Update(ctx, entry)
}

func (s *TimeEntryService) DeleteTimeEntry(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
