package inventory

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/sushanthach12/ecom-go/internal/auth"
	"github.com/sushanthach12/ecom-go/internal/database"
	"github.com/sushanthach12/ecom-go/internal/products"
)

type returnValue struct {
	Adapter Adapter
}

func Register(router *chi.Mux, db *database.DB, logger *slog.Logger, productAdapter products.Adapter) *returnValue {
	repository := newRepository(db)
	adapter := newAdapter(repository, logger)
	service := NewService(logger, repository, productAdapter)
	handler := newHandler(service, logger)

	router.Route("/inventory", func(r chi.Router) {
		r.Use(auth.RequireAuth)

		r.Post("/update", handler.UpdateStock)
	})

	return &returnValue{
		Adapter: adapter,
	}
}
