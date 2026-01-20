package inventory

import (
	"foodlink_backend/middleware"
	"net/http"
	"strings"
)

// SetupRoutes sets up inventory routes
func SetupRoutes(service *Service, handler *Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		switch {
		case path == "/" && r.Method == http.MethodGet:
			handler.GetAll(w, r)
		case path == "/expiring" && r.Method == http.MethodGet:
			handler.GetExpiring(w, r)
		case path == "/expired" && r.Method == http.MethodGet:
			handler.GetExpired(w, r)
		case strings.HasPrefix(path, "/") && len(path) > 1:
			idPath := strings.TrimPrefix(path, "/")
			// Check if it's a UUID (not a special path)
			if len(idPath) == 36 { // UUID length
				switch r.Method {
				case http.MethodGet:
					handler.GetByID(w, r)
				case http.MethodPut:
					handler.Update(w, r)
				case http.MethodDelete:
					handler.Delete(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			} else {
				http.NotFound(w, r)
			}
		case r.Method == http.MethodPost:
			handler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return middleware.Chain(authMiddleware)(mux)
}
