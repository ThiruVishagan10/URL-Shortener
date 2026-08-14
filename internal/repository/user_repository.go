package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error

	FindByGoogleSubject(
		ctx context.Context,
		googleSubject string,
	) (*model.User, error)
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	query := `
		INSERT INTO users (
			google_subject,
			email,
			name,
			avatar_url
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			creaeted_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.GoogleSubject,
		user.Email,
		user.Name,
		user.AvatarURL,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" &&
				pgErr.ConstraintName == "user_google_subject_key" {
				return apperrors.ErrDuplicateUser
			}
		}
		return err
	}
	return nil
}

func (r *PostgresUserRepository) FindByGoogleSubject(
	ctx context.Context,
	googleSubject string,
) (*model.User, error) {
	query := `
		SELECT
			id,
			google_subject,
			email,
			name,
			avatar_url,
			created_at,
			updated_at
		FROM users
		WHERE google_subject = $1
	`

	user := &model.User{}

	err := r.db.QueryRow(
		ctx,
		query,
		googleSubject,
	).Scan(
		&user.ID,
		&user.GoogleSubject,
		&user.Email,
		&user.Name,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}
