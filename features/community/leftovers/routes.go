package leftovers

import (
	"foodlink_backend/middleware"
	"net/http"
	"strings"
)

func SetupRoutes(service *Service, handler *Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/community/leftovers/")
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
		case len(pathParts) == 2 && pathParts[1] == "claim" && r.Method == http.MethodPost:
			handler.CreateClaim(w, r)
		case len(pathParts) == 2 && pathParts[1] == "claims" && r.Method == http.MethodGet:
			handler.GetClaims(w, r)
		case r.Method == http.MethodPost:
			handler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return middleware.Chain(authMiddleware)(mux)
}
