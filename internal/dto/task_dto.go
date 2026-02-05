package dto

import (
	"time"

	"github.com/dmitry/taskmanager/internal/domain"
)

type CreateTaskRequest struct {
	Title        string              `json:"title" validate:"required,min=3,max=500"`
	Description  string              `json:"description"`
	Priority     int                 `json:"priority" validate:"min=0,max=2"`
	DueDate      *string             `json:"due_date,omitempty"`
	Participants []ParticipantInput `json:"participants"`
}

type UpdateTaskRequest struct {
	Title       string  `json:"title" validate:"required,min=3,max=500"`
	Description string  `json:"description"`
	Priority    int     `json:"priority" validate:"min=0,max=2"`
	DueDate     *string `json:"due_date,omitempty"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=new in_progress code_review testing returned_with_errors closed"`
}

type ParticipantInput struct {
	EmployeeID string `json:"employee_id" validate:"required,uuid"`
	Role       string `json:"role" validate:"required,oneof=executor responsible customer"`
}

type AddParticipantRequest struct {
	EmployeeID string `json:"employee_id" validate:"required,uuid"`
	Role       string `json:"role" validate:"required,oneof=executor responsible customer"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	CreatedBy   string    `json:"created_by"`
	Archived    bool      `json:"archived"`
	DueDate     *string   `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToTaskResponse(t *domain.Task) TaskResponse {
	resp := TaskResponse{
		ID:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    t.Priority,
		CreatedBy:   t.CreatedBy.String(),
		Archived:    t.Archived,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	if t.DueDate != nil {
		dueDate := t.DueDate.Format("2006-01-02")
		resp.DueDate = &dueDate
	}

	return resp
}

type TaskParticipantResponse struct {
	ID         string    `json:"id"`
	EmployeeID string    `json:"employee_id"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToTaskParticipantResponse(p *domain.TaskParticipant) TaskParticipantResponse {
	return TaskParticipantResponse{
		ID:         p.ID.String(),
		EmployeeID: p.EmployeeID.String(),
		Role:       string(p.Role),
		CreatedAt:  p.CreatedAt,
	}
}
