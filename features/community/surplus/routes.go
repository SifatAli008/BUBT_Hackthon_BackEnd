package surplus

import (
	"foodlink_backend/middleware"
	"net/http"
	"strings"
)

func SetupRoutes(service *Service, handler *Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/community/surplus/")
		pathParts := strings.Split(path, "/")
		
		switch {
		case path == "" && r.Method == http.MethodGet:
			handler.GetAll(w, r)
		case len(pathParts) == 1 && len(pathParts[0]) == 36 && r.Method == http.MethodGet:
			handler.GetByID(w, r)
		case len(pathParts) == 1 && len(pathParts[0]) == 36 && r.Method == http.MethodPut:
			handler.Update(w, r)
		case len(pathParts) == 1 && len(pathParts[0]) == 36 && r.Method == http.MethodDelete:
			handler.Delete(w, r)
		case len(pathParts) == 2 && pathParts[1] == "request" && r.Method == http.MethodPost:
			handler.CreateRequest(w, r)
		case len(pathParts) == 2 && pathParts[1] == "requests" && r.Method == http.MethodGet:
			handler.GetRequests(w, r)
		case len(pathParts) == 3 && pathParts[1] == "requests" && r.Method == http.MethodPut:
			handler.UpdateRequest(w, r)
		case len(pathParts) == 2 && pathParts[1] == "comments" && r.Method == http.MethodPost:
			handler.CreateComment(w, r)
		case len(pathParts) == 2 && pathParts[1] == "comments" && r.Method == http.MethodGet:
			handler.GetComments(w, r)
		case r.Method == http.MethodPost:
			handler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return middleware.Chain(authMiddleware)(mux)
}
