package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		userID string,
		originalURL string,
	) (*model.URL, error)

	FindByUserID(
		ctx context.Context,
		userID string,
	) ([]*model.URL, error)

	UpdateVisibility(
		ctx context.Context,
		shortID string,
		userID string,
		visibility string,
	) error
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
			original_url,
			user_id,
			visibility
		)
		VALUES ($1, $2, $3, $4)
		RETURNING 
			id,
			created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		url.ShortID,
		url.OriginalURL,
		url.UserID,
		url.Visibility,
	).Scan(
		&url.ID,
		&url.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" &&
				pgErr.ConstraintName == "unique_original_url" {
				return apperrs.ErrDuplicateURL
			}
		}
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
			user_id,
			visibility,
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
		&url.UserID,
		&url.Visibility,
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
	userID string,
	originalURL string,
) (*model.URL, error) {

	query := `
		SELECT 
			id,
			short_id,
			original_url,
			user_id,
			visibility,
			created_at
		FROM urls
		WHERE user_id = $1 
			AND original_url = $2
	`

	url := &model.URL{}

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		originalURL,
	).Scan(
		&url.ID,
		&url.ShortID,
		&url.OriginalURL,
		&url.UserID,
		&url.Visibility,
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

func (r *PostgresURLRepository) FindByUserID(
	ctx context.Context,
	userID string,
) ([]*model.URL, error) {

	query := `
		SELECT
			id,
			short_id,
			original_url,
			user_id,
			visibility,
			created_at
		FROM urls
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*model.URL

	for rows.Next() {
		url := &model.URL{}

		if err := rows.Scan(
			&url.ID,
			&url.ShortID,
			&url.OriginalURL,
			&url.UserID,
			&url.Visibility,
			&url.CreatedAt,
		); err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (r *PostgresURLRepository) UpdateVisibility(
	ctx context.Context,
	shortID string,
	userID string,
	visibility string,
) error {

	query := `
		UPDATE urls
		SET visibility = $1
		WHERE short_id = $2
			AND user_id = $3
	`

	result, err := r.db.Exec(
		ctx,
		query,
		visibility,
		shortID,
		userID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrs.ErrURLNotFound
	}

	return nil
}
