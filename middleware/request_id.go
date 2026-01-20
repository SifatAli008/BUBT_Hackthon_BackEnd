package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// RequestID adds a unique request ID to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(r *http.Request) string {
	if requestID, ok := r.Context().Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
