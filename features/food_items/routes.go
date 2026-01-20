package food_items

import (
	"net/http"
)

// SetupRoutes sets up food items routes
func SetupRoutes(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == "/" {
				handler.GetAll(w, r)
			} else {
				handler.GetByID(w, r)
			}
		case http.MethodPost:
			handler.Create(w, r)
		case http.MethodPut:
			handler.Update(w, r)
		case http.MethodDelete:
			handler.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
