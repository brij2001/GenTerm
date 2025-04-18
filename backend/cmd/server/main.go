package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// Register common MIME types
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".html", "text/html")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".png", "image/png")
	mime.AddExtensionType(".jpg", "image/jpeg")
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".ico", "image/x-icon")

	// Set up API routes with CORS middleware
	http.HandleFunc("/api/chat", api.EnableCors(apiHandler.HandleChat))
	http.HandleFunc("/api/session", api.EnableCors(apiHandler.HandleSession))

	// Create a file server for static files
	staticDir := "/app/frontend/build"
	fs := http.FileServer(http.Dir(staticDir))

	// Handler for static files that also handles SPA routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log request info for debugging
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// If it's an API request, don't handle it here
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Set correct content type for known file extensions
		ext := filepath.Ext(r.URL.Path)
		if ext != "" {
			if mimeType := mime.TypeByExtension(ext); mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
				log.Printf("Set MIME type for %s: %s", r.URL.Path, mimeType)
			}
		}

		// Check if the requested file exists
		path := filepath.Join(staticDir, r.URL.Path)
		stat, err := os.Stat(path)
		if err == nil {
			log.Printf("File found: %s (Size: %d, IsDir: %v)", path, stat.Size(), stat.IsDir())
		} else {
			log.Printf("File not found: %s (Error: %v)", path, err)
		}

		// Special handling for CSS, JS, and other static assets
		if strings.Contains(r.URL.Path, "/static/") || strings.HasSuffix(r.URL.Path, ".js") || strings.HasSuffix(r.URL.Path, ".css") {
			if os.IsNotExist(err) {
				log.Printf("Static file not found: %s", path)
				http.NotFound(w, r)
				return
			}
			log.Printf("Serving static file: %s", path)
			fs.ServeHTTP(w, r)
			return
		}

		// For all other paths, serve index.html (SPA routing)
		if os.IsNotExist(err) || r.URL.Path != "/" && r.URL.Path != "/index.html" {
			indexFile := filepath.Join(staticDir, "index.html")
			log.Printf("Serving index.html for path: %s", r.URL.Path)
			http.ServeFile(w, r, indexFile)
			return
		}

		// Otherwise, serve the requested file
		log.Printf("Serving file directly: %s", path)
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
