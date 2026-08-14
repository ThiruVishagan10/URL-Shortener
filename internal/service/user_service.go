package service

import (
	"context"
	"errors"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
	"github.com/thiruvishagan10/URL-Shortener/internal/repository"
)

type UserService interface {
	FindOrCreate(
		ctx context.Context,
		googleSubject string,
		email string,
		name string,
		avatar_url *string,
	) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(
	repo repository.UserRepository,
) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) FindOrCreate(
	ctx context.Context,
	googleSubject string,
	email string,
	name string,
	avatar_url *string,
) (*model.User, error) {
	existingUser, err := s.repo.FindByGoogleSubject(
		ctx,
		googleSubject,
	)

	if err == nil {
		return existingUser, nil
	}

	if !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, err
	}

	user := &model.User{
		GoogleSubject: googleSubject,
		Email:         email,
		Name:          name,
		AvatarURL:     avatar_url,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, apperrors.ErrDuplicateUser) {
			existingUser, findErr := s.repo.FindByGoogleSubject(
				ctx,
				googleSubject,
			)

			if findErr != nil {
				return nil, findErr
			}

			return existingUser, nil
		}

		return nil, err
	}

	return user, nil
}
