package products

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
)

type returnValue struct {
	Adapter Adapter
}

func Register(router *chi.Mux, db *sql.DB, logger *slog.Logger) *returnValue {
	repository := newRepository(db)
	service := NewService(repository)
	handler := newHandler(service, logger)

	adapter := newAdapter(service)

	router.Get("/products", handler.List)
	router.Get("/products/{product_id}", handler.GetById)

	return &returnValue{
		Adapter: adapter,
	}
}
