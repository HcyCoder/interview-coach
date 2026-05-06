package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// New creates a structured zerolog logger with the requested log level.
//
// Inputs: level is a case-insensitive zerolog level name such as "info" or
// "debug".
// Outputs: a configured zerolog.Logger, or an error when the level string is
// invalid.
func New(level string) (zerolog.Logger, error) {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(strings.TrimSpace(level)))
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("parse log level: %w", err)
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(parsedLevel)

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "backend").
		Logger().
		Level(parsedLevel)

	return logger, nil
}
