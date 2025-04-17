package api

import (
	"encoding/json"
	"net/http"

	"github.com/genterm/backend/internal/config"
	"github.com/genterm/backend/internal/llm"
	"github.com/genterm/backend/internal/session"
)

// Handler manages API endpoints
type Handler struct {
	config         *config.Config
	sessionManager *session.Manager
	llmClient      *llm.Client
}

// MessageContent represents the different types of content in a message
type MessageContent struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

// ImageURL represents an image URL object
type ImageURL struct {
	URL string `json:"url"`
}

// ChatRequest is the structure for chat requests
type ChatRequest struct {
	SessionID      string           `json:"sessionId"`
	Query          string           `json:"query"`
	Context        []string         `json:"context"`
	MessageContent []MessageContent `json:"messageContent,omitempty"`
}

// ChatResponse is the structure for chat responses
type ChatResponse struct {
	SessionID string `json:"sessionId"`
	Response  string `json:"response"`
}

// SessionRequest is the structure for session requests
type SessionRequest struct {
	Action string `json:"action"`
	ID     string `json:"id,omitempty"`
}

// SessionResponse is the structure for session responses
type SessionResponse struct {
	ID       string            `json:"id"`
	Messages []session.Message `json:"messages,omitempty"`
	Error    string            `json:"error,omitempty"`
}

// NewHandler creates a new API handler
func NewHandler(cfg *config.Config, sessionMgr *session.Manager) *Handler {
	return &Handler{
		config:         cfg,
		sessionManager: sessionMgr,
		llmClient:      llm.NewClient(cfg),
	}
}

// HandleChat handles chat requests
func (h *Handler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Ensure we have a valid session
	session, exists := h.sessionManager.GetSession(req.SessionID)
	if !exists {
		http.Error(w, "Invalid session", http.StatusBadRequest)
		return
	}

	// Generate system prompt
	systemPrompt := "You are a helpful assistant. Use the provided context to answer questions accurately."

	// Convert session messages to LLM messages
	var sessionMessages []llm.Message
	for _, msg := range session.Messages {
		sessionMessages = append(sessionMessages, llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	var response string
	var err error

	// Check if request contains image data
	if len(req.MessageContent) > 0 {
		// Add user message with image to session
		h.sessionManager.AddMessage(session.ID, "user", req.Query+" [with image]")

		// Convert MessageContent to ContentItem
		contentItems := make([]llm.ContentItem, len(req.MessageContent))
		for i, content := range req.MessageContent {
			contentItems[i] = llm.ContentItem{
				Type: content.Type,
				Text: content.Text,
				ImageURL: llm.ImageURL{
					URL: content.ImageURL.URL,
				},
			}
		}

		// Get LLM response using multimodal API with conversation history
		response, err = h.llmClient.GenerateMultimodalCompletionWithHistory(sessionMessages, contentItems, req.Context, systemPrompt)
	} else {
		// Add user message to session
		h.sessionManager.AddMessage(session.ID, "user", req.Query)

		// Get LLM response using RAG with conversation history
		response, err = h.llmClient.GenerateCompletionWithHistory(sessionMessages, req.Query, req.Context, systemPrompt)
	}

	if err != nil {
		http.Error(w, "Error generating response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add assistant message to session
	h.sessionManager.AddMessage(session.ID, "assistant", response)

	// Send response
	resp := ChatResponse{
		SessionID: session.ID,
		Response:  response,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleSession handles session management
func (h *Handler) HandleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch req.Action {
	case "create":
		session := h.sessionManager.NewSession()
		json.NewEncoder(w).Encode(SessionResponse{
			ID: session.ID,
		})

	case "get":
		session, exists := h.sessionManager.GetSession(req.ID)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(SessionResponse{
				Error: "Session not found",
			})
			return
		}
		json.NewEncoder(w).Encode(SessionResponse{
			ID:       session.ID,
			Messages: session.Messages,
		})

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}
