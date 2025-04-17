package config

import (
	"errors"
	"os"
)

// Config holds application configuration
type Config struct {
	LLMBaseURL string
	LLMAPIKey  string
	LLMModel   string
	Port       string
}

// NewConfig initializes configuration from environment variables
func NewConfig() (*Config, error) {
	baseURL := os.Getenv("LLM_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1" // Default to OpenAI API
	}

	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		return nil, errors.New("LLM_API_KEY environment variable is required")
	}

	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "gpt-4o" // Default model
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	return &Config{
		LLMBaseURL: baseURL,
		LLMAPIKey:  apiKey,
		LLMModel:   model,
		Port:       port,
	}, nil
}
