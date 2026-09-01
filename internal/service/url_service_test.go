package service

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
	"github.com/thiruvishagan10/URL-Shortener/internal/repository"
)

type mockURLRepository struct {
	deleteFunc func(
		ctx context.Context,
		shortID string,
		userID string,
	) error
}

func (m *mockURLRepository) Create(
	ctx context.Context,
	url *model.URL,
) error {
	return nil
}

func (m *mockURLRepository) FindByShortID(
	ctx context.Context,
	shortID string,
) (*model.URL, error) {
	return nil, nil
}

func (m *mockURLRepository) FindByOriginalURL(
	ctx context.Context,
	userID string,
	originalURL string,
) (*model.URL, error) {
	return nil, nil
}

func (m *mockURLRepository) FindByUserID(
	ctx context.Context,
	userID string,
) ([]*model.URL, error) {
	return nil, nil
}

func (m *mockURLRepository) UpdateVisibility(
	ctx context.Context,
	shortID string,
	userID string,
	visibility string,
) error {
	return nil
}

func (m *mockURLRepository) Delete(
	ctx context.Context,
	shortID string,
	userID string,
) error {
	return m.deleteFunc(ctx, shortID, userID)
}

var _ repository.URLRepository = (*mockURLRepository)(nil)

func TestURLService_Delete(t *testing.T) {
	t.Run("successfully deletes URL owned by user", func(t *testing.T) {
		repo := &mockURLRepository{
			deleteFunc: func(
				ctx context.Context,
				shortID string,
				userID string,
			) error {
				if shortID != "abc123" {
					t.Errorf("expected shortID abc123, got %s", shortID)
				}

				if userID != "user-123" {
					t.Errorf("expected userID user-123, got %s", userID)
				}

				return nil
			},
		}

		service := NewURLService(repo)

		err := service.Delete(
			context.Background(),
			"abc123",
			"user-123",
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns ErrURLNotFound when URL does not belong to user", func(t *testing.T) {
		repo := &mockURLRepository{
			deleteFunc: func(
				ctx context.Context,
				shortID string,
				userID string,
			) error {
				return apperrors.ErrURLNotFound
			},
		}

		service := NewURLService(repo)

		err := service.Delete(
			context.Background(),
			"abc123",
			"user-456",
		)

		if !errors.Is(err, apperrors.ErrURLNotFound) {
			t.Fatalf(
				"expected ErrURLNotFound, got %v",
				err,
			)
		}
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		expectedErr := errors.New("database error")

		repo := &mockURLRepository{
			deleteFunc: func(
				ctx context.Context,
				shortID string,
				userID string,
			) error {
				return expectedErr
			},
		}

		service := NewURLService(repo)

		err := service.Delete(
			context.Background(),
			"abc123",
			"user-123",
		)

		if !errors.Is(err, expectedErr) {
			t.Fatalf(
				"expected database error, got %v",
				err,
			)
		}
	})
}
