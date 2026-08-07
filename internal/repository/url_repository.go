package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrs "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error

	FindByShortID(
		ctx context.Context,
		shortID string,
	) (*model.URL, error)

	FindByOriginalURL(
		ctx context.Context,
		originalURL string,
	) (*model.URL, error)
}

type PostgresURLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *PostgresURLRepository {
	return &PostgresURLRepository{
		db: db,
	}
}

func (r *PostgresURLRepository) Create(
	ctx context.Context,
	url *model.URL,
) error {

	query := `
		INSERT INTO urls (
			short_id,
			original_url
		)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		url.ShortID,
		url.OriginalURL,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresURLRepository) FindByShortID(
	ctx context.Context,
	shortID string,
) (*model.URL, error) {

	query := `
		SELECT 
			id,
			short_id,
			original_url,
			created_at
		FROM urls
		where short_id = $1
	`

	url := &model.URL{}

	err := r.db.QueryRow(
		ctx,
		query,
		shortID,
	).Scan(
		&url.ID,
		&url.ShortID,
		&url.OriginalURL,
		&url.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrs.ErrURLNotFound
		}
		return nil, err
	}
	return url, nil
}

func (r *PostgresURLRepository) FindByOriginalURL(
	ctx context.Context,
	original_url string,
) (*model.URL, error) {

	query := `
		SELECT 
			id,
			short_id,
			original_url,
			created_at
		FROM urls
		WHERE original_url = $1
	`

	url := &model.URL{}

	err := r.db.QueryRow(
		ctx,
		query,
		original_url,
	).Scan(
		&url.ID,
		&url.ShortID,
		&url.OriginalURL,
		&url.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrs.ErrURLNotFound
		}
		return nil, err
	}
	return url, nil
}
