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
	Create(ctx context.Context, originalURL string) (*model.URL, error)

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

	urlModel := &model.URL{
		ShortID:     shortID,
		OriginalURL: originalURL,
	}

	if err := s.repo.Create(ctx, urlModel); err != nil {
		return nil, err
	}

	return urlModel, nil
}

func (s *urlService) GetByShortID(
	ctx context.Context,
	shortID string,
) (*model.URL, error) {
	return s.repo.FindByShortID(ctx, shortID)
}
