package service

import (
	"context"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/internal/repository"
	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/google/uuid"
)

type MessageService struct {
	repo   repository.MessageRepository
	logger *logger.Logger
}

func NewMessageService(repo repository.MessageRepository, logger *logger.Logger) *MessageService {
	return &MessageService{
		repo:   repo,
		logger: logger,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, taskID, authorID uuid.UUID, content string) (*domain.TaskMessage, error) {
	message := domain.NewTaskMessage(taskID, &authorID, content, false)

	if err := s.repo.Create(ctx, message); err != nil {
		return nil, err
	}

	s.logger.Info("Сообщение создано", "message_id", message.ID, "task_id", taskID)

	return message, nil
}

func (s *MessageService) GetTaskMessages(ctx context.Context, taskID uuid.UUID) ([]*domain.TaskMessage, error) {
	return s.repo.GetByTask(ctx, taskID)
}

func (s *MessageService) UpdateMessage(ctx context.Context, messageID uuid.UUID, content string) error {
	message, err := s.repo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	message.Content = content

	return s.repo.Update(ctx, message)
}

func (s *MessageService) DeleteMessage(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
