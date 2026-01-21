package shopping_list

import (
	"foodlink_backend/middleware"
	"net/http"
	"strings"
)

// SetupRoutes sets up shopping list routes
func SetupRoutes(service *Service, handler *Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case path == "/" && r.Method == http.MethodGet:
			handler.GetAll(w, r)
		case path == "/" && r.Method == http.MethodPost:
			handler.Create(w, r)
		case path == "/compute-missing" && r.Method == http.MethodPost:
			handler.ComputeMissing(w, r)
		case strings.HasPrefix(path, "/") && len(path) > 1:
			// /:id or /:id/toggle
			parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
			if len(parts) == 0 {
				http.NotFound(w, r)
				return
			}
			if len(parts) == 2 && parts[1] == "toggle" {
				handler.Toggle(w, r)
				return
			}
			// UUID length check
			if len(parts[0]) == 36 {
				switch r.Method {
				case http.MethodGet:
					handler.GetByID(w, r)
				case http.MethodPut, http.MethodPatch:
					handler.Update(w, r)
				case http.MethodDelete:
					handler.Delete(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}
			http.NotFound(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return middleware.Chain(authMiddleware)(mux)
}

