package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusNew                TaskStatus = "new"
	TaskStatusInProgress         TaskStatus = "in_progress"
	TaskStatusCodeReview         TaskStatus = "code_review"
	TaskStatusTesting            TaskStatus = "testing"
	TaskStatusReturnedWithErrors TaskStatus = "returned_with_errors"
	TaskStatusClosed             TaskStatus = "closed"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusNew, TaskStatusInProgress, TaskStatusCodeReview,
		TaskStatusTesting, TaskStatusReturnedWithErrors, TaskStatusClosed:
		return true
	}
	return false
}

func (s TaskStatus) String() string {
	return string(s)
}

type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Priority    int        `json:"priority"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	Archived    bool       `json:"archived"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewTask(title, description string, priority int, createdBy uuid.UUID, dueDate *time.Time) *Task {
	now := time.Now()
	return &Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      TaskStatusNew,
		Priority:    priority,
		CreatedBy:   createdBy,
		Archived:    false,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
