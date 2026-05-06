package session

import (
	"context"
	"sync"

	"interview-coach/backend/internal/model"
)

// MemoryStore keeps interview sessions in-memory for tests and local development.
type MemoryStore struct {
	mu       sync.RWMutex
	nextID   uint
	sessions map[uint]*model.InterviewSession
}

// NewMemoryStore creates an empty memory-backed session store.
//
// Inputs: none.
// Outputs: a writable session store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nextID:   0,
		sessions: make(map[uint]*model.InterviewSession),
	}
}

// Create stores a new session and returns its generated ID.
//
// Inputs: ctx is unused but preserved for interface compatibility; session is the session record to persist.
// Outputs: the assigned session ID, or an error when persistence fails.
func (s *MemoryStore) Create(_ context.Context, session *model.InterviewSession) (uint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	session.ID = s.nextID
	copy := *session
	s.sessions[session.ID] = &copy
	return session.ID, nil
}

// Get retrieves a session by ID.
//
// Inputs: ctx is unused but preserved for interface compatibility; id is the session identifier.
// Outputs: the session record or ErrSessionNotFound when absent.
func (s *MemoryStore) Get(_ context.Context, id uint) (*model.InterviewSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}

	copy := *session
	return &copy, nil
}

// Update replaces an existing session record.
//
// Inputs: ctx is unused but preserved for interface compatibility; session is the updated record.
// Outputs: nil on success or ErrSessionNotFound when the session does not exist.
func (s *MemoryStore) Update(_ context.Context, session *model.InterviewSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[session.ID]; !ok {
		return ErrSessionNotFound
	}

	copy := *session
	s.sessions[session.ID] = &copy
	return nil
}

var _ SessionStore = (*MemoryStore)(nil)
