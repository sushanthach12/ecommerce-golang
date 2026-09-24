package users

import "context"

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateParams struct {
	Name         string
	Email        string
	PasswordHash string
}

type Service interface {
	Create(ctx context.Context, data CreateParams) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (User, error)
}

type createRepoParams struct {
	Name         string
	Email        string
	PasswordHash string
}

type usersRepository interface {
	FindByEmail(ctx context.Context, email string) (userEntity, error)
	FindByEmailCredentials(ctx context.Context, email string) (userEntity, error)
	Create(ctx context.Context, data createRepoParams) (userEntity, error)
}
