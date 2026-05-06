package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"interview-coach/backend/internal/session"
)

// Client calls the AI service over HTTP.
type Client struct {
	baseURL string
	http    *http.Client
}

// NewClient creates a new AI service client.
//
// Inputs: baseURL is the AI service root URL.
// Outputs: a configured HTTP client for resume parsing and question generation.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ParseResume sends a PDF resume to the AI service and returns the structured profile.
//
// Inputs: ctx is the request context; pdfBytes are the uploaded PDF bytes; filename is the original filename.
// Outputs: a structured resume profile or an error when the downstream call fails.
func (c *Client) ParseResume(
	ctx context.Context,
	pdfBytes []byte,
	filename string,
) (session.ResumeProfile, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return session.ResumeProfile{}, err
	}
	if _, err := part.Write(pdfBytes); err != nil {
		return session.ResumeProfile{}, err
	}
	if err := writer.Close(); err != nil {
		return session.ResumeProfile{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/parse-resume",
		&body,
	)
	if err != nil {
		return session.ResumeProfile{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return session.ResumeProfile{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return session.ResumeProfile{}, fmt.Errorf("ai parse-resume failed: %s", resp.Status)
	}

	var profile aiResumeProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return session.ResumeProfile{}, err
	}

	return profile.toSessionProfile(), nil
}

// GenerateQuestion asks the AI service for the next interview question.
//
// Inputs: ctx is the request context; resumeProfile is the structured resume; conversation is the interview history.
// Outputs: the next interview question or an error when the downstream call fails.
func (c *Client) GenerateQuestion(
	ctx context.Context,
	resumeProfile session.ResumeProfile,
	conversation []session.ConversationTurn,
) (string, error) {
	payload := map[string]any{
		"resume_profile": resumeProfile,
		"conversation":   conversation,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/generate-question",
		buf,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("ai generate-question failed: %s", resp.Status)
	}

	var out struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.Question, nil
}

type aiResumeProfile struct {
	TechStack           []string          `json:"tech_stack"`
	Projects            []aiResumeProject `json:"projects"`
	Summary             string            `json:"summary"`
	Strengths           []string          `json:"strengths"`
	SuggestedFocusAreas []string          `json:"suggested_focus_areas"`
}

type aiResumeProject struct {
	Name       string   `json:"name"`
	Summary    string   `json:"summary"`
	TechStack  []string `json:"tech_stack"`
	Highlights []string `json:"highlights"`
}

func (profile aiResumeProfile) toSessionProfile() session.ResumeProfile {
	projects := make([]session.ResumeProject, 0, len(profile.Projects))
	for _, project := range profile.Projects {
		projects = append(projects, session.ResumeProject{
			Name:       project.Name,
			Summary:    project.Summary,
			TechStack:  project.TechStack,
			Highlights: project.Highlights,
		})
	}

	return session.ResumeProfile{
		TechStack:           profile.TechStack,
		Projects:            projects,
		Summary:             profile.Summary,
		Strengths:           profile.Strengths,
		SuggestedFocusAreas: profile.SuggestedFocusAreas,
	}
}

var _ session.AIClient = (*Client)(nil)
