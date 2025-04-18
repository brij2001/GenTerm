package main

import (
	"log"
	"net/http"
	"os"

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

	// Set up HTTP server with CORS middleware
	http.HandleFunc("/api/chat", api.EnableCors(apiHandler.HandleChat))
	http.HandleFunc("/api/session", api.EnableCors(apiHandler.HandleSession))

	// Create a file server for static files
	fs := http.FileServer(http.Dir("../../frontend/dist"))

	// Serve static files for the frontend
	http.Handle("/", fs)

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
