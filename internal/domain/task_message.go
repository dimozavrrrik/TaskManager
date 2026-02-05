package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskMessage struct {
	ID              uuid.UUID  `json:"id"`
	TaskID          uuid.UUID  `json:"task_id"`
	AuthorID        *uuid.UUID `json:"author_id,omitempty"`
	Content         string     `json:"content"`
	IsSystemMessage bool       `json:"is_system_message"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewTaskMessage(taskID uuid.UUID, authorID *uuid.UUID, content string, isSystemMessage bool) *TaskMessage {
	now := time.Now()
	return &TaskMessage{
		ID:              uuid.New(),
		TaskID:          taskID,
		AuthorID:        authorID,
		Content:         content,
		IsSystemMessage: isSystemMessage,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func NewSystemMessage(taskID uuid.UUID, content string) *TaskMessage {
	return NewTaskMessage(taskID, nil, content, true)
}
