package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"

	"github.com/thiruvishagan10/URL-Shortener/internal/config"
)

type GoogleOAuthProvider struct {
	config *oauth2.Config
}

type AuthorizationRequest struct {
	URL          string
	State        string
	CodeVerifier string
}

type GoogleIdentity struct {
	Subject   string
	Email     string
	Name      string
	AvatarURL *string
}

func NewGoogleOAuthProvider(cfg *config.Config) *GoogleOAuthProvider {
	return &GoogleOAuthProvider{
		config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"openid",
				"profile",
				"email",
			},
		},
	}
}

func (p *GoogleOAuthProvider) Begin() (*AuthorizationRequest, error) {
	state, err := generateState()
	if err != nil {
		return nil, err
	}

	codeVerifier := oauth2.GenerateVerifier()

	url := p.config.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(codeVerifier),
	)

	return &AuthorizationRequest{
		URL:          url,
		State:        state,
		CodeVerifier: codeVerifier,
	}, nil
}

func generateState() (string, error) {
	buffer := make([]byte, 32)

	if _, err := rand.Read(buffer); err != nil {
		return "", errors.New("failed tp genereate OAuth state")
	}

	return base64.RawStdEncoding.EncodeToString(buffer), nil
}

func (p *GoogleOAuthProvider) Authenticate(
	ctx context.Context,
	code string,
	codeVerifier string,
) (*GoogleIdentity, error) {
	token, err := p.config.Exchange(
		ctx,
		code,
		oauth2.VerifierOption(codeVerifier),
	)
	if err != nil {
		return nil, errors.New("failed to exchange Google authorizaton code")
	}

	if !token.Valid() {
		return nil, errors.New("invalid Google token")
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, errors.New("Google ID token missing")
	}

	payload, err := idtoken.Validate(
		ctx,
		rawIDToken,
		p.config.ClientID,
	)
	if err != nil {
		return nil, errors.New("invalid Google ID token")
	}

	if payload.Subject == "" {
		return nil, errors.New("Google subject missing")
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	var avatarURL *string

	if picture != "" {
		avatarURL = &picture
	}

	return &GoogleIdentity{
		Subject:   payload.Subject,
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
	}, nil
}
