package session

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"interview-coach/backend/internal/model"
)

// GormStore persists sessions through GORM.
type GormStore struct {
	db *gorm.DB
}

// NewGormStore wraps a GORM database with the session store interface.
//
// Inputs: db is an initialized GORM connection.
// Outputs: a durable session store.
func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

// Create stores a new session row in MySQL.
//
// Inputs: ctx is the request context; session is the record to persist.
// Outputs: the generated session ID, or an error when persistence fails.
func (s *GormStore) Create(ctx context.Context, session *model.InterviewSession) (uint, error) {
	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return 0, err
	}

	return session.ID, nil
}

// Get retrieves a session row from MySQL.
//
// Inputs: ctx is the request context; id is the session identifier.
// Outputs: the session record or ErrSessionNotFound when absent.
func (s *GormStore) Get(ctx context.Context, id uint) (*model.InterviewSession, error) {
	var record model.InterviewSession
	if err := s.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Update persists changes to an existing session row.
//
// Inputs: ctx is the request context; session is the updated record.
// Outputs: nil on success or an error when persistence fails.
func (s *GormStore) Update(ctx context.Context, session *model.InterviewSession) error {
	return s.db.WithContext(ctx).Save(session).Error
}

var _ SessionStore = (*GormStore)(nil)
