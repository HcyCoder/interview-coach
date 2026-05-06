package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/rs/zerolog"
)

// Health handles GET /health and returns the service status payload.
//
// Inputs: ctx is the request context from Hertz; c is the current request
// context.
// Outputs: a 200 JSON response with {"status":"ok"}.
func Health(ctx context.Context, c *app.RequestContext) {
	_ = ctx
	c.JSON(http.StatusOK, struct {
		Status string `json:"status"`
	}{Status: "ok"})
}

// NotFoundHandler logs unknown routes and returns a JSON 404 response.
type NotFoundHandler struct {
	Logger zerolog.Logger
}

// NewNotFoundHandler creates a route-not-found handler with structured logging.
//
// Inputs: logger is the service logger used to record 404 paths.
// Outputs: a NotFoundHandler ready for Hertz NoRoute registration.
func NewNotFoundHandler(logger zerolog.Logger) NotFoundHandler {
	return NotFoundHandler{Logger: logger}
}

// Handle records the missing route and returns a 404 JSON error payload.
//
// Inputs: ctx is the request context from Hertz; c is the current request
// context.
// Outputs: a 404 JSON response with {"error":"route not found"}.
func (h NotFoundHandler) Handle(ctx context.Context, c *app.RequestContext) {
	_ = ctx

	h.Logger.Warn().
		Str("method", string(c.Method())).
		Str("path", string(c.Path())).
		Msg("route not found")

	c.JSON(http.StatusNotFound, struct {
		Error string `json:"error"`
	}{Error: "route not found"})
}
