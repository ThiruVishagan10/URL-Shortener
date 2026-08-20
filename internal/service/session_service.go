package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
	"github.com/thiruvishagan10/URL-Shortener/internal/repository"
)

const sessionLifetime = 24 * time.Hour

type SessionService interface {
	Create(
		ctx context.Context,
		userID string,
	) (*model.Session, error)

	Find(
		ctx context.Context,
		sessionID string,
	) (*model.Session, error)

	Delete(
		ctx context.Context,
		sessionID string,
	) error
}

type sessionService struct {
	repo repository.SessionRepository
}

func NewSessionService(
	repo repository.SessionRepository,
) SessionService {
	return &sessionService{
		repo: repo,
	}
}

func (s *sessionService) Create(
	ctx context.Context,
	userID string,
) (*model.Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &model.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionLifetime),
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionService) Find(
	ctx context.Context,
	sessionID string,
) (*model.Session, error) {
	session, err := s.repo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.repo.Delete(ctx, sessionID)

		return nil, apperrors.ErrSessionNotFound
	}

	return session, nil
}

func (s *sessionService) Delete(
	ctx context.Context,
	sessionID string,
) error {
	return s.repo.Delete(ctx, sessionID)
}

func generateSessionID() (string, error) {
	buffer := make([]byte, 32)

	if _, err := rand.Read(buffer); err != nil {
		return "", errors.New("failed to generate session ID")
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
