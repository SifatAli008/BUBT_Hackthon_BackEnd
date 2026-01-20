package auth

import (
	"foodlink_backend/middleware"
	"net/http"
)

// SetupRoutes sets up authentication routes
func SetupRoutes(service *Service, handler *Handler) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/login", handler.Login)

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/logout", handler.Logout)
	protectedMux.HandleFunc("/refresh", handler.RefreshToken)
	protectedMux.HandleFunc("/me", handler.GetMe)

	// Apply auth middleware to protected routes
	protectedHandler := middleware.Chain(
		AuthMiddleware(service),
	)(protectedMux)

	// Mount protected routes
	mux.Handle("/", protectedHandler)

	return mux
}
