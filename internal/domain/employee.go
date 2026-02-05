package domain

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Department   string     `json:"department"`
	Position     string     `json:"position"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Никогда не выводить в JSON
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewEmployee(name, department, position, email string) *Employee {
	now := time.Now()
	return &Employee{
		ID:         uuid.New(),
		Name:       name,
		Department: department,
		Position:   position,
		Email:      email,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
