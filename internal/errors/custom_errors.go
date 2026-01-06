package errors

import (
	"be/internal/constants"
	"fmt"
)

type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewNotFound(msg string) *AppError {
	return &AppError{
		Code:    constants.ErrNotFound,
		Message: msg,
		Status:  404,
	}
}

func NewBadRequest(msg string) *AppError {
	return &AppError{
		Code:    constants.ErrBadRequest,
		Message: msg,
		Status:  400,
	}
}

func NewAuthorized(msg string) *AppError {
	return &AppError{
		Code:    constants.ErrUnauthorized,
		Message: msg,
		Status:  401,
	}
}

func NewVerificationRequired(msg string) *AppError {
	return &AppError{
		Code:    constants.ErrVerificationRequired,
		Message: msg,
		Status:  400,
	}
}

func NewInternalServerError(msg string) *AppError {
	return &AppError{
		Code:    constants.ErrInternalServerError,
		Message: msg,
		Status:  500,
	}
}

func NewForbiden(msg string) *AppError {
	return &AppError{
		Code:    constants.ErrForbiden,
		Message: msg,
		Status:  403,
	}
}

func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status
	}
	return 500
}
