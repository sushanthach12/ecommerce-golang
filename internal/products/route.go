package products

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/sushanthach12/ecom-go/internal/database"
)

type returnValue struct {
	Adapter Adapter
}

func Register(router *chi.Mux, db *database.DB, logger *slog.Logger) *returnValue {
	repository := newRepository(db)
	service := NewService(repository, logger)
	handler := newHandler(service, logger)

	adapter := newAdapter(repository, logger)

	router.Get("/products", handler.List)
	router.Get("/products/{product_id}", handler.GetById)

	return &returnValue{
		Adapter: adapter,
	}
}
