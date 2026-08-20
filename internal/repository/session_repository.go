package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error

	FindByID(
		ctx context.Context,
		sessionID string,
	) (*model.Session, error)

	Delete(
		ctx context.Context,
		sessionID string,
	) error
}

type PostgresSessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *PostgresSessionRepository {
	return &PostgresSessionRepository{
		db: db,
	}
}

func (r *PostgresSessionRepository) Create(
	ctx context.Context,
	session *model.Session,
) error {
	query := `
		INSERT INTO sessions (
			user_id,
			expires_at
		)
		VALUES ($1, $2)
		RETURNING 
			id,
			created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		session.UserID,
		session.ExpiresAt,
	).Scan(
		&session.ID,
		&session.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresSessionRepository) FindByID(
	ctx context.Context,
	sessionID string,
) (*model.Session, error) {
	query := `
	SELECT 
		id,
		user_id,
		expires_at,
		created_at
	FROM sessions
	WHERE id = $1
	`

	session := &model.Session{}

	err := r.db.QueryRow(
		ctx,
		query,
		sessionID,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrSessionNotFound
		}
		return nil, err
	}
	return session, nil
}

func (r *PostgresSessionRepository) Delete(
	ctx context.Context,
	sessionID string,
) error {
	query := `
		DELETE FROM sessions
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		sessionID,
	)

	return err
}
