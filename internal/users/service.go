package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type service struct {
	repo   usersRepository
	logger *slog.Logger
}

func newService(repo usersRepository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Create(ctx context.Context, payload CreateParams) (User, error) {

	// find by email
	_, err := s.repo.FindByEmail(ctx, payload.Email)

	switch {
	case err == nil:
		// no error means a row was found -> email already taken
		s.logger.Warn("service: create: email already registered", "email", payload.Email)
		return User{}, ErrUserAlreadyExists

	case errors.Is(err, ErrUserNotFound):
		// expected — email is free, fall through and continue creating

	default:
		s.logger.Error("service: create: failed to check existing email", "error", err)
		return User{}, fmt.Errorf("something went wrong")
	}

	user, err := s.repo.Create(ctx, createRepoParams{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: payload.PasswordHash,
	})
	if err != nil {
		s.logger.Error("service: failed to create", "error", err)
		return User{}, fmt.Errorf("something went wrong")
	}

	return User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (User, error) {
	entity, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return User{}, ErrUserNotFound
		}
		s.logger.Error("service: getByEmail: failed", "error", err)
		return User{}, fmt.Errorf("something went wrong")
	}

	return User{
		ID:        entity.ID,
		Name:      entity.Name,
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

// GetCredentialsByEmail returns the user including PasswordHash.
// Intended only for the auth package's login flow — do not use this
// for anything that returns data to a client.
func (s *service) GetCredentialsByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.FindByEmailCredentials(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return User{}, ErrUserNotFound
		}
		s.logger.Error("service: getCredentialsByEmail: failed", "error", err)
		return User{}, fmt.Errorf("something went wrong")
	}

	return User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}
