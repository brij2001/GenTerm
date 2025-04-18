package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/genterm/backend/internal/api"
	"github.com/genterm/backend/internal/config"
	"github.com/genterm/backend/internal/session"
	"github.com/joho/godotenv"
)

// These variables will be set during build via ldflags
var (
	LlmApiKey  string
	LlmModel   string
	LlmBaseUrl string
)

func main() {
	// Load environment variables from .env if available
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Set environment variables from build-time values if not already set
	if os.Getenv("LLM_API_KEY") == "" && LlmApiKey != "" {
		os.Setenv("LLM_API_KEY", LlmApiKey)
	}

	if os.Getenv("LLM_MODEL") == "" && LlmModel != "" {
		os.Setenv("LLM_MODEL", LlmModel)
	}

	if os.Getenv("LLM_BASE_URL") == "" && LlmBaseUrl != "" {
		os.Setenv("LLM_BASE_URL", LlmBaseUrl)
	}

	// Initialize configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize session manager
	sessionManager := session.NewManager()

	// Initialize API handlers
	apiHandler := api.NewHandler(cfg, sessionManager)

	// Set up API routes with CORS middleware
	http.HandleFunc("/api/chat", api.EnableCors(apiHandler.HandleChat))
	http.HandleFunc("/api/session", api.EnableCors(apiHandler.HandleSession))

	// Create a file server for static files
	staticDir := "/app/frontend/build"
	fs := http.FileServer(http.Dir(staticDir))

	// Handler for static files that also handles SPA routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the requested file exists
		path := filepath.Join(staticDir, r.URL.Path)
		_, err := os.Stat(path)

		// If path doesn't exist or is a directory, serve the index.html
		if os.IsNotExist(err) || r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		// Otherwise, serve the requested file
		fs.ServeHTTP(w, r)
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	// Start the server
	log.Printf("Starting server on :%s", port)
	log.Printf("MODEL: %s", model)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
