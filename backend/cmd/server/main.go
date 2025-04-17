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

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
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

	// Set up HTTP server
	http.HandleFunc("/api/chat", apiHandler.HandleChat)
	http.HandleFunc("/api/session", apiHandler.HandleSession)

	// Serve static files for the frontend
	fs := http.FileServer(http.Dir("../../frontend/dist"))
	http.Handle("/", fs)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting server on :%s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
