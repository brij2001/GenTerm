package api

import (
	"net/http"
	"os"
)

// CorsMiddleware wraps an http.Handler to add CORS headers
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check origin and set appropriate CORS headers
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"https://genterm-d9eue4hhc7azcvb4.eastus-01.azurewebsites.net",
			"http://localhost:3000",
		}

		// Allow any origin in development mode or check against allowed list
		if os.Getenv("NODE_ENV") != "production" || containsOrigin(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0])
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Helper function to check if origin is in allowed list
func containsOrigin(allowedOrigins []string, origin string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}

// EnableCors is a wrapper for individual handlers
func EnableCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check origin and set appropriate CORS headers
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"https://genterm-d9eue4hhc7azcvb4.eastus-01.azurewebsites.net",
			"http://localhost:3000",
		}

		// Allow any origin in development mode or check against allowed list
		if os.Getenv("NODE_ENV") != "production" || containsOrigin(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0])
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}
