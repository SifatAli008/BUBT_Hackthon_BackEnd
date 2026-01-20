package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorResponse sends an error JSON response
func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Message: message,
		Error:   err,
	}

	json.NewEncoder(w).Encode(response)
}

// JSONResponse sends a JSON response
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// CreatedResponse sends a 201 Created response
func CreatedResponse(w http.ResponseWriter, message string, data interface{}) {
	SuccessResponse(w, http.StatusCreated, message, data)
}

// OKResponse sends a 200 OK response
func OKResponse(w http.ResponseWriter, message string, data interface{}) {
	SuccessResponse(w, http.StatusOK, message, data)
}

// NoContentResponse sends a 204 No Content response
func NoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequestResponse sends a 400 Bad Request response
func BadRequestResponse(w http.ResponseWriter, message string, err interface{}) {
	ErrorResponse(w, http.StatusBadRequest, message, err)
}

// UnauthorizedResponse sends a 401 Unauthorized response
func UnauthorizedResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusUnauthorized, message, nil)
}

// ForbiddenResponse sends a 403 Forbidden response
func ForbiddenResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusForbidden, message, nil)
}

// NotFoundResponse sends a 404 Not Found response
func NotFoundResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusNotFound, message, nil)
}

// ConflictResponse sends a 409 Conflict response
func ConflictResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusConflict, message, nil)
}

// InternalServerErrorResponse sends a 500 Internal Server Error response
func InternalServerErrorResponse(w http.ResponseWriter, message string, err interface{}) {
	ErrorResponse(w, http.StatusInternalServerError, message, err)
}
