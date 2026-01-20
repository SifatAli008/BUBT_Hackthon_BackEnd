package auth

import (
	"context"
	"foodlink_backend/errors"
	"foodlink_backend/utils"
	"net/http"
	"strings"
)

// AuthMiddleware validates JWT tokens and sets user in context
func AuthMiddleware(service *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.UnauthorizedResponse(w, "Authorization header required")
				return
			}

			// Check Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.UnauthorizedResponse(w, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Validate token and get user
			user, err := service.ValidateToken(tokenString)
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok {
					utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
					return
				}
				utils.UnauthorizedResponse(w, "Invalid token")
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*User)
			if !ok {
				utils.UnauthorizedResponse(w, "Authentication required")
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range requiredRoles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				utils.ForbiddenResponse(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuth middleware validates token if present but doesn't require it
func OptionalAuth(service *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					user, err := service.ValidateToken(parts[1])
					if err == nil {
						ctx := context.WithValue(r.Context(), "user", user)
						r = r.WithContext(ctx)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
