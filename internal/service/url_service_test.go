package service

import (
	"context"
	"testing"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
)

type fakeURLRepository struct {
	findByOriginalURL func(
		ctx context.Context,
		originalURL string,
	) (*model.URL, error)

	create func(
		ctx context.Context,
		url *model.URL,
	) error
}

func (f *fakeURLRepository) Create(
	ctx context.Context,
	url *model.URL,
) error {
	return f.create(ctx, url)
}

func (f *fakeURLRepository) FindByShortID(
	ctx context.Context,
	shortID string,
) (*model.URL, error) {
	return nil, nil
}

func (f *fakeURLRepository) FindByOriginalURL(
	ctx context.Context,
	originalURL string,
) (*model.URL, error) {
	return f.findByOriginalURL(ctx, originalURL)
}

func TestURLService_Create(t *testing.T) {

	repo := &fakeURLRepository{
		findByOriginalURL: func(
			ctx context.Context,
			originalURL string,
		) (*model.URL, error) {
			return nil, apperrors.ErrURLNotFound
		},

		create: func(
			ctx context.Context,
			url *model.URL,
		) error {
			return nil
		},
	}

	service := NewURLService(repo)

	result, err := service.Create(
		context.Background(),
		"https://example.com",
	)

	if err != nil {
		t.Fatalf("expected no more, got %v", err)
	}

	if result == nil {
		t.Fatal("expected URL, got nil")
	}

	if result.OriginalURL != "https://example.com" {
		t.Fatalf(
			"expected original url %q, got %q",
			"https://example.com",
			result.OriginalURL,
		)
	}

	if result.ShortID == "" {
		t.Fatal("expected ShortID to be generated")
	}
}
