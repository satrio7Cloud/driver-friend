package errors

import (
	"be/internal/constants"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

func NewForbidden(msg string) *AppError {
	return &AppError{
		Code:    constants.ErrForbiden,
		Message: msg,
		Status:  403,
	}
}

func HandleError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		ctx.JSON(appErr.Status, gin.H{
			"error": appErr.Message,
		})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
	})
}

func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status
	}
	return 500
}
