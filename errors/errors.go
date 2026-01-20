package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWithErr creates a new application error with underlying error
func NewAppErrorWithErr(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
var (
	// 400 Bad Request
	ErrBadRequest          = NewAppError(http.StatusBadRequest, "Bad request")
	ErrInvalidInput        = NewAppError(http.StatusBadRequest, "Invalid input")
	ErrValidationFailed   = NewAppError(http.StatusBadRequest, "Validation failed")
	ErrInvalidJSON         = NewAppError(http.StatusBadRequest, "Invalid JSON")
	ErrMissingRequiredField = NewAppError(http.StatusBadRequest, "Missing required field")

	// 401 Unauthorized
	ErrUnauthorized       = NewAppError(http.StatusUnauthorized, "Unauthorized")
	ErrInvalidToken       = NewAppError(http.StatusUnauthorized, "Invalid token")
	ErrTokenExpired       = NewAppError(http.StatusUnauthorized, "Token expired")
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "Invalid credentials")

	// 403 Forbidden
	ErrForbidden        = NewAppError(http.StatusForbidden, "Forbidden")
	ErrInsufficientPerms = NewAppError(http.StatusForbidden, "Insufficient permissions")

	// 404 Not Found
	ErrNotFound      = NewAppError(http.StatusNotFound, "Resource not found")
	ErrUserNotFound  = NewAppError(http.StatusNotFound, "User not found")
	ErrRecordNotFound = NewAppError(http.StatusNotFound, "Record not found")

	// 409 Conflict
	ErrConflict     = NewAppError(http.StatusConflict, "Resource conflict")
	ErrAlreadyExists = NewAppError(http.StatusConflict, "Resource already exists")
	ErrDuplicateKey  = NewAppError(http.StatusConflict, "Duplicate key")

	// 500 Internal Server Error
	ErrInternalServer = NewAppError(http.StatusInternalServerError, "Internal server error")
	ErrDatabase       = NewAppError(http.StatusInternalServerError, "Database error")
	ErrUnexpected     = NewAppError(http.StatusInternalServerError, "Unexpected error")
)

// WrapError wraps an error with an AppError
func WrapError(err error, appErr *AppError) *AppError {
	if err == nil {
		return appErr
	}
	return NewAppErrorWithErr(appErr.Code, appErr.Message, err)
}
