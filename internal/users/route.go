package users

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/sushanthach12/ecom-go/internal/database"
)

type returnValue struct {
	Adapter Service
}

func Register(router *chi.Mux, db *database.DB, logger *slog.Logger) *returnValue {
	repository := newRepository(db)
	service := newService(repository, logger)
	handler := newHandler(service, logger)

	router.Get("/users/me", handler.Me)

	return &returnValue{
		Adapter: service,
	}
}
