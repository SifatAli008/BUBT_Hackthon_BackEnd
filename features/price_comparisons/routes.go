package price_comparisons

import (
	"net/http"
	"strings"
)

func SetupRoutes(handler *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/price-comparisons/")
		switch {
		case path == "" && r.Method == http.MethodGet:
			handler.GetAll(w, r)
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
	return mux
}
