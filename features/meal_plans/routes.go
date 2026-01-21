package meal_plans

import (
	"foodlink_backend/middleware"
	"net/http"
	"strings"
)

func SetupRoutes(service *Service, handler *Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case path == "/weekly" && r.Method == http.MethodGet:
			handler.GetWeekly(w, r)
		case path == "/" && r.Method == http.MethodPost:
			handler.Upsert(w, r)
		case strings.HasPrefix(path, "/") && len(path) > 1:
			idPath := strings.TrimPrefix(path, "/")
			if len(idPath) == 36 && r.Method == http.MethodDelete {
				handler.Delete(w, r)
				return
			}
			http.NotFound(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return middleware.Chain(authMiddleware)(mux)
}

