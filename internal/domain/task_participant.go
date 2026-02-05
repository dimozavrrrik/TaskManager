package domain

import (
	"time"

	"github.com/google/uuid"
)

type ParticipantRole string

const (
	ParticipantRoleExecutor    ParticipantRole = "executor"
	ParticipantRoleResponsible ParticipantRole = "responsible"
	ParticipantRoleCustomer    ParticipantRole = "customer"
)

func (r ParticipantRole) IsValid() bool {
	switch r {
	case ParticipantRoleExecutor, ParticipantRoleResponsible, ParticipantRoleCustomer:
		return true
	}
	return false
}

func (r ParticipantRole) String() string {
	return string(r)
}

type TaskParticipant struct {
	ID         uuid.UUID       `json:"id"`
	TaskID     uuid.UUID       `json:"task_id"`
	EmployeeID uuid.UUID       `json:"employee_id"`
	Role       ParticipantRole `json:"role"`
	CreatedAt  time.Time       `json:"created_at"`
}

func NewTaskParticipant(taskID, employeeID uuid.UUID, role ParticipantRole) *TaskParticipant {
	return &TaskParticipant{
		ID:         uuid.New(),
		TaskID:     taskID,
		EmployeeID: employeeID,
		Role:       role,
		CreatedAt:  time.Now(),
	}
}
