package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL        string
	Port               string
	FRONTEND_URL       string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	frontend_url := os.Getenv("FRONTEND_URL")

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	if port == "" {
		return nil, errors.New("PORT is required")
	}

	if frontend_url == "" {
		return nil, errors.New("FRONTEND_URL is required")
	}

	if googleClientID == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID is required")
	}

	if googleClientSecret == "" {
		return nil, errors.New("GOOGLE_CLIENT_SECRET is required")
	}

	if googleRedirectURL == "" {
		return nil, errors.New("GOOGLE_REDIRECT_URL is required")
	}

	return &Config{
		DatabaseURL:        dbURL,
		Port:               port,
		FRONTEND_URL:       frontend_url,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectURL:  googleRedirectURL,
	}, nil
}
