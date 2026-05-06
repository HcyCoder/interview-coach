package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config stores runtime settings loaded from environment variables.
type Config struct {
	ServerAddr        string
	LogLevel          string
	DBDSN             string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

// Load reads environment variables and returns the parsed runtime configuration.
//
// Inputs: none.
// Outputs: a Config populated from environment variables, or an error when a
// required value is missing or a numeric/duration value is malformed.
func Load() (Config, error) {
	cfg := Config{
		ServerAddr:        getEnvOrDefault("SERVER_ADDR", ":8080"),
		LogLevel:          getEnvOrDefault("LOG_LEVEL", "info"),
		DBMaxOpenConns:    10,
		DBMaxIdleConns:    10,
		DBConnMaxLifetime: 30 * time.Minute,
	}

	cfg.DBDSN = os.Getenv("DB_DSN")
	if cfg.DBDSN == "" {
		return Config{}, errors.New("DB_DSN is required")
	}

	if value := os.Getenv("DB_MAX_OPEN_CONNS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, fmt.Errorf("parse DB_MAX_OPEN_CONNS: %w", err)
		}
		cfg.DBMaxOpenConns = parsed
	}

	if value := os.Getenv("DB_MAX_IDLE_CONNS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, fmt.Errorf("parse DB_MAX_IDLE_CONNS: %w", err)
		}
		cfg.DBMaxIdleConns = parsed
	}

	if value := os.Getenv("DB_CONN_MAX_LIFETIME"); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("parse DB_CONN_MAX_LIFETIME: %w", err)
		}
		cfg.DBConnMaxLifetime = parsed
	}

	return cfg, nil
}

// getEnvOrDefault returns the value of a variable when it exists, otherwise it
// returns the provided default value.
//
// Inputs: name is the environment variable key; fallback is the default value.
// Outputs: the environment value when set, or fallback when unset.
func getEnvOrDefault(name string, fallback string) string {
	value, ok := os.LookupEnv(name)
	if !ok || value == "" {
		return fallback
	}

	return value
}
