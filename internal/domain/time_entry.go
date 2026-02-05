package domain

import (
	"time"

	"github.com/google/uuid"
)

type TimeEntry struct {
	ID          uuid.UUID  `json:"id"`
	TaskID      uuid.UUID  `json:"task_id"`
	EmployeeID  uuid.UUID  `json:"employee_id"`
	Hours       float64    `json:"hours"`
	Description string     `json:"description"`
	EntryDate   time.Time  `json:"entry_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewTimeEntry(taskID, employeeID uuid.UUID, hours float64, description string, entryDate time.Time) *TimeEntry {
	now := time.Now()
	return &TimeEntry{
		ID:          uuid.New(),
		TaskID:      taskID,
		EmployeeID:  employeeID,
		Hours:       hours,
		Description: description,
		EntryDate:   entryDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type TimeSummary struct {
	TaskID          uuid.UUID `json:"task_id"`
	TotalHours      float64   `json:"total_hours"`
	EntryCount      int       `json:"entry_count"`
	UniqueEmployees int       `json:"unique_employees"`
}
