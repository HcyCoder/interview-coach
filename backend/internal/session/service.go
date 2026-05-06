package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"interview-coach/backend/internal/model"
)

var (
	// ErrSessionNotFound indicates the requested session does not exist.
	ErrSessionNotFound = errors.New("session not found")
	// ErrMessageRequired indicates the user message is empty.
	ErrMessageRequired = errors.New("message is required")
)

// Service coordinates AI calls and session persistence.
type Service struct {
	ai      AIClient
	store   SessionStore
	options Options
	now     func() time.Time
}

// NewService creates a session service.
//
// Inputs: store persists sessions; ai calls the AI service; options carry defaults.
// Outputs: a ready-to-use session service.
func NewService(store SessionStore, ai AIClient, options Options) *Service {
	return &Service{
		ai:      ai,
		store:   store,
		options: options,
		now:     time.Now,
	}
}

// CreateSession parses a resume, stores the session, and returns the first question.
//
// Inputs: ctx is the request context; filename is the uploaded PDF name; pdfBytes are the file contents.
// Outputs: a structured result or an error.
func (s *Service) CreateSession(
	ctx context.Context,
	filename string,
	pdfBytes []byte,
) (CreateSessionResult, error) {
	resumeProfile, err := s.ai.ParseResume(ctx, pdfBytes, filename)
	if err != nil {
		return CreateSessionResult{}, err
	}

	firstQuestion, err := s.ai.GenerateQuestion(ctx, resumeProfile, nil)
	if err != nil {
		return CreateSessionResult{}, err
	}

	conversation := []ConversationTurn{
		{Role: "assistant", Content: firstQuestion},
	}

	resumeJSON, err := json.Marshal(resumeProfile)
	if err != nil {
		return CreateSessionResult{}, fmt.Errorf("marshal resume profile: %w", err)
	}

	conversationJSON, err := json.Marshal(conversation)
	if err != nil {
		return CreateSessionResult{}, fmt.Errorf("marshal conversation: %w", err)
	}

	startedAt := s.now()
	sessionRecord := &model.InterviewSession{
		UserID:           s.options.DefaultUserID,
		ResumeJSON:       resumeJSON,
		ConversationJSON: conversationJSON,
		CurrentQuestion:  firstQuestion,
		Status:           "ACTIVE",
		StartedAt:        &startedAt,
	}

	sessionID, err := s.store.Create(ctx, sessionRecord)
	if err != nil {
		return CreateSessionResult{}, err
	}

	return CreateSessionResult{
		SessionID:     sessionID,
		ResumeProfile: resumeProfile,
		Question:      firstQuestion,
	}, nil
}

// ChatSession appends a candidate reply and returns the next interview question.
//
// Inputs: ctx is the request context; sessionID identifies the interview session; message is the user's answer.
// Outputs: a structured chat result or an error.
func (s *Service) ChatSession(
	ctx context.Context,
	sessionID uint,
	message string,
) (ChatSessionResult, error) {
	if strings.TrimSpace(message) == "" {
		return ChatSessionResult{}, ErrMessageRequired
	}

	sessionRecord, err := s.store.Get(ctx, sessionID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return ChatSessionResult{}, ErrSessionNotFound
		}
		return ChatSessionResult{}, err
	}

	var resumeProfile ResumeProfile
	if err := json.Unmarshal(sessionRecord.ResumeJSON, &resumeProfile); err != nil {
		return ChatSessionResult{}, fmt.Errorf("unmarshal resume profile: %w", err)
	}

	var conversation []ConversationTurn
	if len(sessionRecord.ConversationJSON) > 0 {
		if err := json.Unmarshal(sessionRecord.ConversationJSON, &conversation); err != nil {
			return ChatSessionResult{}, fmt.Errorf("unmarshal conversation: %w", err)
		}
	}

	conversation = append(conversation, ConversationTurn{Role: "user", Content: message})
	nextQuestion, err := s.ai.GenerateQuestion(ctx, resumeProfile, conversation)
	if err != nil {
		return ChatSessionResult{}, err
	}

	conversation = append(conversation, ConversationTurn{Role: "assistant", Content: nextQuestion})

	conversationJSON, err := json.Marshal(conversation)
	if err != nil {
		return ChatSessionResult{}, fmt.Errorf("marshal conversation: %w", err)
	}

	sessionRecord.ConversationJSON = conversationJSON
	sessionRecord.CurrentQuestion = nextQuestion
	sessionRecord.UpdatedAt = s.now()

	if err := s.store.Update(ctx, sessionRecord); err != nil {
		return ChatSessionResult{}, err
	}

	return ChatSessionResult{
		SessionID: sessionID,
		Question:  nextQuestion,
	}, nil
}
