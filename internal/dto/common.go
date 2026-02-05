package dto

import "github.com/dmitry/taskmanager/pkg/errors"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool               `json:"success"`
	Error   ErrorDetailWrapper `json:"error"`
}

type ErrorDetailWrapper struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Details []errors.ErrorDetail `json:"details,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
