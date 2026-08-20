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
	) (*model.URL, error)

	GetByShortID(
		ctx context.Context,
		shortID string,
	) (*model.URL, error)
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
) (*model.URL, error) {

	existingURL, err := s.repo.FindByOriginalURL(
		ctx,
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
	}

	if err := s.repo.Create(ctx, url); err != nil {

		if errors.Is(err, apperrors.ErrDuplicateURL) {
			existingURL, findErr := s.repo.FindByOriginalURL(
				ctx,
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
) (*model.URL, error) {
	return s.repo.FindByShortID(ctx, shortID)
}
