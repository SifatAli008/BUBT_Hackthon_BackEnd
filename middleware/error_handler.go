package middleware

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/utils"
	"log"
	"net/http"
)

// ErrorHandler handles errors and returns appropriate HTTP responses
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to capture errors
		rw := &errorResponseWriter{
			ResponseWriter: w,
			statusCode:    http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		// If an error was set, handle it
		if rw.err != nil {
			handleError(w, r, rw.err, rw.statusCode)
		}
	})
}

type errorResponseWriter struct {
	http.ResponseWriter
	statusCode int
	err        error
}

func (rw *errorResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *errorResponseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

// handleError processes and sends error responses
func handleError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	requestID := GetRequestID(r)

	// Log the error
	log.Printf("Error [%s]: %v", requestID, err)

	// Check if it's an AppError
	if appErr, ok := err.(*errors.AppError); ok {
		errorResponse := map[string]interface{}{
			"code":      appErr.Code,
			"message":   appErr.Message,
			"request_id": requestID,
		}

		if appErr.Err != nil {
			errorResponse["details"] = appErr.Err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":    false,
			"error":      errorResponse,
			"request_id": requestID,
		})
		return
	}

	// Handle generic errors
	errorResponse := map[string]interface{}{
		"code":       statusCode,
		"message":    err.Error(),
		"request_id": requestID,
	}

	// Don't expose internal errors in production
	if statusCode >= http.StatusInternalServerError {
		errorResponse["message"] = "Internal server error"
		log.Printf("Internal error [%s]: %v", requestID, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    false,
		"error":      errorResponse,
		"request_id": requestID,
	})
}

// RecoverPanic recovers from panics and returns error responses
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(r)
				log.Printf("Panic recovered [%s]: %v", requestID, err)

				utils.InternalServerErrorResponse(
					w,
					"Internal server error",
					map[string]interface{}{
						"request_id": requestID,
					},
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
