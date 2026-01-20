package middleware

import (
	"net/http"
)

// CORS handles Cross-Origin Resource Sharing
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CORSWithConfig handles CORS with custom configuration
func CORSWithConfig(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			if len(allowedOrigins) == 0 {
				allowed = true // Allow all if no origins specified
			} else {
				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin || allowedOrigin == "*" {
						allowed = true
						break
					}
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// Set allowed methods
			methods := "GET, POST, PUT, DELETE, PATCH, OPTIONS"
			if len(allowedMethods) > 0 {
				methods = ""
				for i, method := range allowedMethods {
					if i > 0 {
						methods += ", "
					}
					methods += method
				}
			}
			w.Header().Set("Access-Control-Allow-Methods", methods)

			// Set allowed headers
			headers := "Content-Type, Authorization, X-Request-ID"
			if len(allowedHeaders) > 0 {
				headers = ""
				for i, header := range allowedHeaders {
					if i > 0 {
						headers += ", "
					}
					headers += header
				}
			}
			w.Header().Set("Access-Control-Allow-Headers", headers)
			w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
