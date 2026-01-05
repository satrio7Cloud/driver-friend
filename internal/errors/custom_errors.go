package errors

import "fmt"

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
		Code:    "NOT_FOUND",
		Message: msg,
		Status:  404,
	}
}

func NewBadRequest(msg string) *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: msg,
		Status:  400,
	}
}

func NewAuthorized(msg string) *AppError {
	return &AppError{
		Code:    "UNAUTHORIZED",
		Message: msg,
		Status:  401,
	}
}

func NewVerificationRequired(msg string) *AppError {
	return &AppError{
		Code:    "BAD_REQEUST",
		Message: msg,
		Status:  400,
	}
}

func NewInternalServerError(msg string) *AppError {
	return &AppError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: msg,
		Status:  500,
	}
}

func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status
	}
	return 500
}
