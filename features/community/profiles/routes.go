package profiles

import (
	"foodlink_backend/middleware"
	"net/http"
	"strings"
)

func SetupRoutes(service *Service, handler *Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/community/profile/")
		
		switch {
		case path == "" && r.Method == http.MethodGet:
			handler.GetProfile(w, r)
		case path == "" && r.Method == http.MethodPost:
			handler.CreateProfile(w, r)
		case path == "" && r.Method == http.MethodPut:
			handler.UpdateProfile(w, r)
		case path != "" && r.Method == http.MethodGet:
			handler.GetProfileByUsername(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return middleware.Chain(authMiddleware)(mux)
}
