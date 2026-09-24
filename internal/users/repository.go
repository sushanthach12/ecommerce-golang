package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sushanthach12/ecom-go/internal/database"
)

type repository struct {
	db *database.DB
}

func newRepository(db *database.DB) usersRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (userEntity, error) {
	u := userEntity{}

	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, created_at FROM users WHERE email = $1`, email,
	)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userEntity{}, ErrUserNotFound
		}
		return userEntity{}, fmt.Errorf("repository: findByEmail: scan failed: %w", err)
	}

	return u, nil
}

func (r *repository) FindByEmailCredentials(ctx context.Context, email string) (userEntity, error) {
	u := userEntity{}

	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`, email,
	)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userEntity{}, ErrUserNotFound
		}
		return userEntity{}, fmt.Errorf("repository: findByEmail: scan failed: %w", err)
	}

	return u, nil
}

func (r *repository) Create(ctx context.Context, data createRepoParams) (userEntity, error) {
	u := userEntity{}

	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)
		 RETURNING id, name, email, created_at, updated_at`,
		data.Name, data.Email, data.PasswordHash,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}
