package dto

import (
	"time"

	"github.com/dmitry/taskmanager/internal/domain"
)

type CreateEmployeeRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Department string `json:"department" validate:"required,max=100"`
	Position   string `json:"position" validate:"required,max=100"`
	Email      string `json:"email" validate:"required,email"`
}

type UpdateEmployeeRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Department string `json:"department" validate:"required,max=100"`
	Position   string `json:"position" validate:"required,max=100"`
	Email      string `json:"email" validate:"required,email"`
}

type EmployeeResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToEmployeeResponse(e *domain.Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:         e.ID.String(),
		Name:       e.Name,
		Department: e.Department,
		Position:   e.Position,
		Email:      e.Email,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
