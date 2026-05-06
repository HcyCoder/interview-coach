package session

import (
	"context"

	"interview-coach/backend/internal/model"
)

// ResumeProject describes a project extracted from a candidate resume.
type ResumeProject struct {
	Name       string   `json:"name"`
	Summary    string   `json:"summary"`
	TechStack  []string `json:"techStack,omitempty"`
	Highlights []string `json:"highlights,omitempty"`
}

// ResumeProfile is the structured resume payload shared across the app.
type ResumeProfile struct {
	TechStack           []string        `json:"techStack"`
	Projects            []ResumeProject `json:"projects"`
	Summary             string          `json:"summary"`
	Strengths           []string        `json:"strengths"`
	SuggestedFocusAreas []string        `json:"suggestedFocusAreas"`
}

// ConversationTurn represents one turn in an interview chat.
type ConversationTurn struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIClient abstracts the AI service integration used by the backend.
type AIClient interface {
	ParseResume(ctx context.Context, pdfBytes []byte, filename string) (ResumeProfile, error)
	GenerateQuestion(
		ctx context.Context,
		resumeProfile ResumeProfile,
		conversation []ConversationTurn,
	) (string, error)
}

// SessionStore persists interview sessions.
type SessionStore interface {
	Create(ctx context.Context, session *model.InterviewSession) (uint, error)
	Get(ctx context.Context, id uint) (*model.InterviewSession, error)
	Update(ctx context.Context, session *model.InterviewSession) error
}

// Options configures session service defaults.
type Options struct {
	DefaultUserID uint
}

// CreateSessionResult is returned after a successful session bootstrap.
type CreateSessionResult struct {
	SessionID     uint
	ResumeProfile ResumeProfile
	Question      string
}

// ChatSessionResult is returned after generating the next interview question.
type ChatSessionResult struct {
	SessionID uint
	Question  string
}
