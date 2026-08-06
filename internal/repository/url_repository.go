package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thiruvishagan10/URL-Shortener/internal/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
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
