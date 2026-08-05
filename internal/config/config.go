package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	if port == "" {
		return nil, errors.New("PORT is required")
	}

	return &Config{
		DatabaseURL: dbURL,
		Port:        port,
	}, nil
}