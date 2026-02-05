package dto

import (
	"time"

	"github.com/dmitry/taskmanager/internal/domain"
)

type CreateMessageRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type UpdateMessageRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type MessageResponse struct {
	ID              string    `json:"id"`
	TaskID          string    `json:"task_id"`
	AuthorID        *string   `json:"author_id,omitempty"`
	Content         string    `json:"content"`
	IsSystemMessage bool      `json:"is_system_message"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToMessageResponse(m *domain.TaskMessage) MessageResponse {
	resp := MessageResponse{
		ID:              m.ID.String(),
		TaskID:          m.TaskID.String(),
		Content:         m.Content,
		IsSystemMessage: m.IsSystemMessage,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}

	if m.AuthorID != nil {
		authorID := m.AuthorID.String()
		resp.AuthorID = &authorID
	}

	return resp
}
