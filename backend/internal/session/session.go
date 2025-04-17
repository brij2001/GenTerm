package session

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Message represents a single chat message
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Session represents a user session with conversation history
type Session struct {
	ID        string    `json:"id"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Manager handles session creation and retrieval
type Manager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
}

// NewManager creates a new session manager
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

// NewSession creates a new session
func (m *Manager) NewSession() *Session {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	sessionID := uuid.New().String()
	now := time.Now()

	session := &Session{
		ID:        sessionID,
		Messages:  []Message{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	m.sessions[sessionID] = session
	return session
}

// GetSession retrieves a session by ID
func (m *Manager) GetSession(id string) (*Session, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.sessions[id]
	return session, exists
}

// AddMessage adds a message to a session
func (m *Manager) AddMessage(sessionID string, role, content string) (*Message, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, false
	}

	now := time.Now()
	message := Message{
		Role:      role,
		Content:   content,
		Timestamp: now,
	}

	session.Messages = append(session.Messages, message)
	session.UpdatedAt = now

	return &message, true
}

// GetMessages retrieves all messages for a session
func (m *Manager) GetMessages(sessionID string) ([]Message, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, false
	}

	return session.Messages, true
}
