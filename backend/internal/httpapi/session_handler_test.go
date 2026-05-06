package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"interview-coach/backend/internal/httpapi"
	"interview-coach/backend/internal/model"
	"interview-coach/backend/internal/session"
)

func TestCreateSessionHandler(t *testing.T) {
	store := session.NewMemoryStore()
	client := &fakeAIClient{
		profile: session.ResumeProfile{
			TechStack: []string{"Go", "MySQL", "Redis"},
			Projects: []session.ResumeProject{
				{
					Name:      "Interview Coach",
					Summary:   "Built an AI interview practice platform",
					TechStack: []string{"Go", "FastAPI"},
				},
			},
			Summary: "Strong backend and platform experience.",
		},
		question: "Walk me through Interview Coach and the trade-offs you made while working with Go, MySQL, Redis?",
	}
	service := session.NewService(store, client, session.Options{DefaultUserID: 1})
	handler := httpapi.NewSessionHandler(service)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("resume", "resume.pdf")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := io.WriteString(part, "%PDF-1.4 mock"); err != nil {
		t.Fatalf("write file bytes: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/sessions", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.CreateSession(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var payload struct {
		SessionID     uint   `json:"sessionId"`
		Question      string `json:"question"`
		ResumeProfile struct {
			TechStack []string `json:"techStack"`
		} `json:"resumeProfile"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.SessionID != 1 {
		t.Fatalf("session id = %d, want 1", payload.SessionID)
	}
	if payload.Question == "" {
		t.Fatal("question is empty")
	}
	if got := strings.Join(payload.ResumeProfile.TechStack, ","); got != "Go,MySQL,Redis" {
		t.Fatalf("tech stack = %q, want Go,MySQL,Redis", got)
	}
}

func TestChatSessionHandler(t *testing.T) {
	store := session.NewMemoryStore()
	client := &fakeAIClient{
		question: "What trade-offs did you make, and how would you scale this design if traffic doubled?",
	}
	service := session.NewService(store, client, session.Options{DefaultUserID: 1})
	handler := httpapi.NewSessionHandler(service)

	_, err := store.Create(context.Background(), &model.InterviewSession{
		UserID:           1,
		ResumeJSON:       []byte(`{"techStack":["Go","MySQL"],"projects":[{"name":"Interview Coach","summary":"Built an AI interview practice platform","techStack":["Go"]}],"summary":"Strong backend and platform experience."}`),
		Status:           "ACTIVE",
		ConversationJSON: []byte(`[{"role":"assistant","content":"Walk me through Interview Coach and the trade-offs you made while working with Go, MySQL?"}]`),
		CurrentQuestion:  "Walk me through Interview Coach and the trade-offs you made while working with Go, MySQL?",
	})
	if err != nil {
		t.Fatalf("seed session: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/sessions/1/chat",
		strings.NewReader(`{"message":"I built the backend and tuned MySQL indexes."}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ChatSession(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Fatalf("content-type = %q, want text/plain", ct)
	}
	if body := rec.Body.String(); !strings.Contains(body, "scale this design") {
		t.Fatalf("stream body = %q, want follow-up question", body)
	}

	updated, err := store.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if !strings.Contains(string(updated.ConversationJSON), "I built the backend") {
		t.Fatalf("conversation was not updated: %s", string(updated.ConversationJSON))
	}
}

type fakeAIClient struct {
	profile  session.ResumeProfile
	question string
}

func (f *fakeAIClient) ParseResume(_ context.Context, _ []byte, _ string) (session.ResumeProfile, error) {
	return f.profile, nil
}

func (f *fakeAIClient) GenerateQuestion(
	_ context.Context,
	_ session.ResumeProfile,
	_ []session.ConversationTurn,
) (string, error) {
	return f.question, nil
}
