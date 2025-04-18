package api

import (
	"log"
	"net/http"
)

// CorsMiddleware wraps an http.Handler to add CORS headers
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the origin
		origin := r.Header.Get("Origin")
		log.Printf("Request origin: %s, Method: %s", origin, r.Method)

		// Set CORS headers - in production, allow same-site requests and specific domains
		if origin != "" {
			// Allow the specific origin that made the request
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin") // Important when using specific origin values
		} else {
			// Fallback for no origin header
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			log.Printf("Handling OPTIONS preflight request")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// EnableCors is a wrapper for individual handlers
func EnableCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the origin
		origin := r.Header.Get("Origin")
		log.Printf("API request origin: %s, Method: %s, Path: %s", origin, r.Method, r.URL.Path)

		// Set CORS headers - in production, allow same-site requests and specific domains
		if origin != "" {
			// Allow the specific origin that made the request
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin") // Important when using specific origin values
		} else {
			// Fallback for no origin header
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			log.Printf("Handling OPTIONS preflight request for API")
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}
