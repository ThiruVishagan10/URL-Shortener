package service

import (
	"context"
	"errors"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/generator"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
	"github.com/thiruvishagan10/URL-Shortener/internal/repository"
)

type URLService interface {
	Create(
		ctx context.Context,
		userID string,
		originalURL string,
		visibility string,
	) (*model.URL, error)

	GetByShortID(
		ctx context.Context,
		shortID string,
		requesterID string,
	) (*model.URL, error)

	GetByUserID(
		ctx context.Context,
		userID string,
	) ([]*model.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) Create(
	ctx context.Context,
	userID string,
	originalURL string,
	visibility string,
) (*model.URL, error) {

	existingURL, err := s.repo.FindByOriginalURL(
		ctx,
		userID,
		originalURL,
	)

	if err == nil {
		return existingURL, nil
	}

	if !errors.Is(err, apperrors.ErrURLNotFound) {
		return nil, err
	}

	shortID, err := generator.Generate(generator.DefaultShortIDLength)
	if err != nil {
		return nil, err
	}

	url := &model.URL{
		ShortID:     shortID,
		OriginalURL: originalURL,
		UserID:      userID,
		Visibility:  visibility,
	}

	if err := s.repo.Create(ctx, url); err != nil {

		if errors.Is(err, apperrors.ErrDuplicateURL) {
			existingURL, findErr := s.repo.FindByOriginalURL(
				ctx,
				userID,
				originalURL,
			)

			if findErr != nil {
				return nil, findErr
			}

			return existingURL, nil
		}
		return nil, err
	}

	return url, nil
}

func (s *urlService) GetByShortID(
	ctx context.Context,
	shortID string,
	requesterID string,
) (*model.URL, error) {
	url, err := s.repo.FindByShortID(ctx, shortID)
	if err != nil {
		return nil, err
	}

	if url.Visibility == model.URLVisibilityPublic {
		return url, nil
	}

	if url.Visibility == model.URLVisibilityPrivate &&
		url.UserID == requesterID {
		return url, err
	}

	return nil, apperrors.ErrURLNotFound
}

func (s *urlService) GetByUserID(
	ctx context.Context,
	userID string,
) ([]*model.URL, error) {
	urls, err := s.repo.FindByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return urls, nil
}
