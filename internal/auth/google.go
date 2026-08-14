package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

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
