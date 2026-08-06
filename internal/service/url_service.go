package service

import (
	"context"

	"github.com/thiruvishagan10/URL-Shortener/internal/generator"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
	"github.com/thiruvishagan10/URL-Shortener/internal/repository"
)

type URLService interface {
	Create(ctx context.Context, originalURL string) (*model.URL, error)
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
