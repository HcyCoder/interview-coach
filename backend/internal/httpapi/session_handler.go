package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"interview-coach/backend/internal/session"
)

// SessionHandler exposes the interview session API surface.
type SessionHandler struct {
	service *session.Service
}

// NewSessionHandler builds a session HTTP handler wrapper.
//
// Inputs: service coordinates AI calls and persistence.
// Outputs: a handler that can be mounted under Hertz or tested with net/http.
func NewSessionHandler(service *session.Service) *SessionHandler {
	return &SessionHandler{service: service}
}

// CreateSession handles resume upload and returns the first interview question.
//
// Inputs: w is the HTTP response writer; r is the incoming multipart request.
// Outputs: a JSON response with the created session and first question.
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("resume")
	if err != nil {
		writeError(w, http.StatusBadRequest, "resume file is required")
		return
	}
	defer file.Close()

	if !isPDFUpload(header.Header.Get("Content-Type"), header.Filename) {
		writeError(w, http.StatusUnsupportedMediaType, "resume must be a PDF")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read resume file")
		return
	}

	result, err := h.service.CreateSession(r.Context(), header.Filename, fileBytes)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"sessionId":     result.SessionID,
		"resumeProfile": result.ResumeProfile,
		"question":      result.Question,
	})
}

// ChatSession handles an answer and streams the next interview question.
//
// Inputs: w is the HTTP response writer; r is the incoming JSON request.
// Outputs: a streamed plain-text response with the next question.
func (h *SessionHandler) ChatSession(w http.ResponseWriter, r *http.Request) {
	sessionID, err := parseSessionID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	var payload struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	result, err := h.service.ChatSession(r.Context(), sessionID, payload.Message)
	if err != nil {
		switch {
		case errors.Is(err, session.ErrSessionNotFound):
			writeError(w, http.StatusNotFound, "session not found")
		case errors.Is(err, session.ErrMessageRequired):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusBadGateway, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	flusher, _ := w.(http.Flusher)
	for _, chunk := range chunkText(result.Question, 18) {
		_, _ = io.WriteString(w, chunk)
		if flusher != nil {
			flusher.Flush()
		}
	}
}

func parseSessionID(path string) (uint, error) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) != 4 || segments[0] != "api" || segments[1] != "sessions" || segments[3] != "chat" {
		return 0, fmt.Errorf("invalid chat path")
	}

	value, err := strconv.ParseUint(segments[2], 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(value), nil
}

func isPDFUpload(contentType, filename string) bool {
	if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		return true
	}

	switch contentType {
	case "application/pdf", "application/x-pdf", "application/octet-stream", "":
		return true
	default:
		return false
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func chunkText(text string, width int) []string {
	if text == "" {
		return []string{""}
	}

	if width <= 0 {
		return []string{text}
	}

	chunks := make([]string, 0, (len(text)/width)+1)
	for len(text) > 0 {
		if utf8.RuneCountInString(text) <= width {
			chunks = append(chunks, text)
			break
		}

		cut := 0
		count := 0
		for i := range text {
			if count == width {
				cut = i
				break
			}
			count++
		}
		if cut == 0 {
			chunks = append(chunks, text)
			break
		}
		chunks = append(chunks, text[:cut])
		text = text[cut:]
	}

	return chunks
}
