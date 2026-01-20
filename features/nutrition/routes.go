package nutrition

import (
	"foodlink_backend/middleware"
	"net/http"
	"strings"
)

func SetupRoutes(service *Service, handler *Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/nutrition/")
		switch {
		case path == "" && r.Method == http.MethodGet:
			handler.GetAll(w, r)
		case path == "today" && r.Method == http.MethodGet:
			handler.GetToday(w, r)
		case path == "stats" && r.Method == http.MethodGet:
			handler.GetStats(w, r)
		case len(path) == 36 && r.Method == http.MethodGet:
			handler.GetByID(w, r)
		case len(path) == 36 && r.Method == http.MethodPut:
			handler.Update(w, r)
		case r.Method == http.MethodPost:
			handler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return middleware.Chain(authMiddleware)(mux)
}
