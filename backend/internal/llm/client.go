package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/genterm/backend/internal/config"
)

// Client for interacting with LLM APIs
type Client struct {
	config *config.Config
	client *http.Client
}

// Message represents a single message in the conversation
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ContentItem represents a single content item in a message
type ContentItem struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

// ImageURL represents an image URL object
type ImageURL struct {
	URL string `json:"url"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
}

// Choice represents a response choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// NewClient creates a new LLM client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		client: &http.Client{},
	}
}

// GenerateCompletion generates a chat completion response
func (c *Client) GenerateCompletion(messages []Message) (string, error) {
	chatRequest := ChatRequest{
		Model:     c.config.LLMModel,
		Messages:  messages,
		MaxTokens: 2000,
	}

	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %w", err)
	}

	// Debug: print request
	// fmt.Printf("Debug - API Request: %s\n", string(jsonData))

	url := fmt.Sprintf("%s/chat/completions", c.config.LLMBaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.LLMAPIKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s, status code: %d", string(bodyBytes), resp.StatusCode)
	}

	var chatResponse ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices returned in response")
	}

	content := chatResponse.Choices[0].Message.Content
	if strContent, ok := content.(string); ok {
		return strContent, nil
	}

	// Try to marshal the content if it's not a string
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("error marshalling content: %w", err)
	}

	return string(contentBytes), nil
}

// GenerateCompletionWithHistory generates a chat completion response using conversation history
func (c *Client) GenerateCompletionWithHistory(sessionMessages []Message, query string, context []string, systemPrompt string) (string, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// Add context as user messages if provided
	if len(context) > 0 {
		contextMessage := "Context information:\n\n"
		for i, ctx := range context {
			contextMessage += fmt.Sprintf("[%d] %s\n\n", i+1, ctx)
		}
		messages = append(messages, Message{
			Role:    "user",
			Content: contextMessage,
		})
	}

	// Add conversation history
	messages = append(messages, sessionMessages...)

	// Add user query as the latest message
	messages = append(messages, Message{
		Role:    "user",
		Content: query,
	})

	return c.GenerateCompletion(messages)
}

// GenerateRAGCompletion generates a completion with RAG context
func (c *Client) GenerateRAGCompletion(query string, context []string, systemPrompt string) (string, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// Add context as user messages
	contextMessage := "Context information:\n\n"
	for i, ctx := range context {
		contextMessage += fmt.Sprintf("[%d] %s\n\n", i+1, ctx)
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: contextMessage,
	})

	// Add user query
	messages = append(messages, Message{
		Role:    "user",
		Content: query,
	})

	return c.GenerateCompletion(messages)
}

// GenerateMultimodalCompletion generates a completion with image and text
func (c *Client) GenerateMultimodalCompletion(messageContent []ContentItem, context []string, systemPrompt string) (string, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// Add context as user messages if provided
	if len(context) > 0 {
		contextMessage := "Context information:\n\n"
		for i, ctx := range context {
			contextMessage += fmt.Sprintf("[%d] %s\n\n", i+1, ctx)
		}
		messages = append(messages, Message{
			Role:    "user",
			Content: contextMessage,
		})
	}

	// Add multimodal content as a single message
	messages = append(messages, Message{
		Role:    "user",
		Content: messageContent,
	})

	return c.GenerateCompletion(messages)
}

// GenerateMultimodalCompletionWithHistory generates a completion with image, text and conversation history
func (c *Client) GenerateMultimodalCompletionWithHistory(sessionMessages []Message, messageContent []ContentItem, context []string, systemPrompt string) (string, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// Add context as user messages if provided
	if len(context) > 0 {
		contextMessage := "Context information:\n\n"
		for i, ctx := range context {
			contextMessage += fmt.Sprintf("[%d] %s\n\n", i+1, ctx)
		}
		messages = append(messages, Message{
			Role:    "user",
			Content: contextMessage,
		})
	}

	// Add conversation history
	messages = append(messages, sessionMessages...)

	// Add multimodal content as a single message
	messages = append(messages, Message{
		Role:    "user",
		Content: messageContent,
	})

	return c.GenerateCompletion(messages)
}
