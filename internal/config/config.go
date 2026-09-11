package config

import (
	"errors"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL        string
	Port               string
	FRONTEND_URL       string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	AppBaseURL         string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	frontend_url := os.Getenv("FRONTEND_URL")
	appBaseURL := strings.TrimRight(os.Getenv("APP_BASE_URL"), "/")

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

	parsedAppBaseURL, err := url.ParseRequestURI(appBaseURL)
	if appBaseURL == "" || err != nil ||
		(parsedAppBaseURL.Scheme != "http" && parsedAppBaseURL.Scheme != "https") ||
		parsedAppBaseURL.Host == "" {
		return nil, errors.New("APP_BASE_URL must be an absolute http or https URL")
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
		AppBaseURL:         appBaseURL,
	}, nil
}
