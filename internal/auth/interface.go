package auth

import (
	"context"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

type registerParams struct {
	Name     string
	Email    string
	Password string
}

type loginParams struct {
	Email    string
	Password string
}

type authUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type authResponse struct {
	User         authUser `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

type Service interface {
	Register(ctx context.Context, payload registerParams) (constants.SimpleResponse[authResponse], error)
	Login(ctx context.Context, payload loginParams) (constants.SimpleResponse[authResponse], error)
}
