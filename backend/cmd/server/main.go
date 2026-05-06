package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/adaptor"

	"interview-coach/backend/internal/ai"
	"interview-coach/backend/internal/config"
	"interview-coach/backend/internal/db"
	"interview-coach/backend/internal/handler"
	"interview-coach/backend/internal/httpapi"
	"interview-coach/backend/internal/logging"
	"interview-coach/backend/internal/session"
)

// main bootstraps the Hertz backend service.
//
// Inputs: none.
// Outputs: none. The process exits with a non-zero code when startup fails.
func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := logging.New(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create logger: %v\n", err)
		os.Exit(1)
	}

	database, err := db.Open(cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("open mysql connection failed")
		os.Exit(1)
	}

	if err := db.AutoMigrate(database); err != nil {
		logger.Error().Err(err).Msg("auto migrate failed")
		os.Exit(1)
	}

	aiClient := ai.NewClient(cfg.AIServiceBaseURL)
	sessionStore := session.NewGormStore(database)
	sessionService := session.NewService(sessionStore, aiClient, session.Options{
		DefaultUserID: cfg.DefaultUserID,
	})
	sessionHandler := httpapi.NewSessionHandler(sessionService)

	h := server.Default(server.WithHostPorts(cfg.ServerAddr))
	h.GET("/health", handler.Health)
	h.NoRoute(handler.NewNotFoundHandler(logger).Handle)
	h.POST("/api/sessions", adaptor.HertzHandler(http.HandlerFunc(sessionHandler.CreateSession)))
	h.POST("/api/sessions/:id/chat", adaptor.HertzHandler(http.HandlerFunc(sessionHandler.ChatSession)))

	logger.Info().Str("addr", cfg.ServerAddr).Msg("backend server starting")

	h.Spin()
}
