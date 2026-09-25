package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/httpx"
	"github.com/sushanthach12/ecom-go/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userService users.Service
	logger      *slog.Logger
}

func newService(userService users.Service, logger *slog.Logger) Service {
	return &service{
		userService: userService,
		logger:      logger,
	}
}

func (s *service) Register(ctx context.Context, payload registerParams) (constants.SimpleResponse[authResponse], error) {
	var defaultResponse constants.SimpleResponse[authResponse]

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return defaultResponse, fmt.Errorf("Something went wrong")
	}

	user, err := s.userService.Create(ctx, users.CreateParams{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: string(hash),
	})
	if err != nil {
		s.logger.Error("auth_service: Failed to create", "error", err)
		if err == users.ErrUserAlreadyExists {
			return defaultResponse, &httpx.ConflictError{
				Field:   "email",
				Message: "user with this email already exists",
			}
		}
		return defaultResponse, err
	}

	access, _ := GenerateAccessToken(user.ID, user.Email)
	refresh, _ := GenerateRefreshToken(user.ID)

	return constants.NewResponse(authResponse{
		User: authUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil), nil
}

func (s *service) Login(ctx context.Context, payload loginParams) (constants.SimpleResponse[authResponse], error) {
	var defaultResponse constants.SimpleResponse[authResponse]

	exists, err := s.userService.GetCredentialsByEmail(ctx, payload.Email)
	if err != nil {
		if err == users.ErrUserNotFound {
			return defaultResponse, &httpx.NotFoundError{
				Message: "please register instead!",
			}
		}
		return defaultResponse, err
	}

	decryptedErr := bcrypt.CompareHashAndPassword([]byte(exists.PasswordHash), []byte(payload.Password))
	if decryptedErr != nil {
		s.logger.Error("failed to decrypt", "error", decryptedErr)
		return defaultResponse, &httpx.UnauthorizedError{
			Message: "invalid credentials",
		}
	}

	access, _ := GenerateAccessToken(exists.ID, exists.Email)
	refresh, _ := GenerateRefreshToken(exists.ID)

	return constants.NewResponse(authResponse{
		User: authUser{
			ID:    exists.ID,
			Name:  exists.Name,
			Email: exists.Email,
		},
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil), nil
}
