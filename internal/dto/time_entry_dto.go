package dto

import (
	"time"

	"github.com/dmitry/taskmanager/internal/domain"
)

type CreateTimeEntryRequest struct {
	Hours       float64 `json:"hours" validate:"required,gt=0"`
	Description string  `json:"description"`
	EntryDate   string  `json:"entry_date" validate:"required"`
}

type UpdateTimeEntryRequest struct {
	Hours       float64 `json:"hours" validate:"required,gt=0"`
	Description string  `json:"description"`
	EntryDate   string  `json:"entry_date" validate:"required"`
}

type TimeEntryResponse struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"task_id"`
	EmployeeID  string    `json:"employee_id"`
	Hours       float64   `json:"hours"`
	Description string    `json:"description"`
	EntryDate   string    `json:"entry_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToTimeEntryResponse(e *domain.TimeEntry) TimeEntryResponse {
	return TimeEntryResponse{
		ID:          e.ID.String(),
		TaskID:      e.TaskID.String(),
		EmployeeID:  e.EmployeeID.String(),
		Hours:       e.Hours,
		Description: e.Description,
		EntryDate:   e.EntryDate.Format("2006-01-02"),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

type TimeSummaryResponse struct {
	TaskID          string  `json:"task_id"`
	TotalHours      float64 `json:"total_hours"`
	EntryCount      int     `json:"entry_count"`
	UniqueEmployees int     `json:"unique_employees"`
}

func ToTimeSummaryResponse(s *domain.TimeSummary) TimeSummaryResponse {
	return TimeSummaryResponse{
		TaskID:          s.TaskID.String(),
		TotalHours:      s.TotalHours,
		EntryCount:      s.EntryCount,
		UniqueEmployees: s.UniqueEmployees,
	}
}
