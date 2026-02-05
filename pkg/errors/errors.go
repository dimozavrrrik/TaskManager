package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
)

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Code    ErrorCode     `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
	Err     error         `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) HTTPStatusCode() int {
	switch e.Code {
	case ErrCodeValidation, ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: message,
	}
}

func Validation(message string, details []ErrorDetail) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
		Details: details,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
	}
}

func Internal(err error, message string) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
		Err:     err,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:    ErrCodeBadRequest,
		Message: message,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    ErrCodeUnauthorized,
		Message: message,
	}
}
