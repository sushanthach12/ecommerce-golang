package orders

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/sushanthach12/ecom-go/internal/auth"
	"github.com/sushanthach12/ecom-go/internal/database"
	"github.com/sushanthach12/ecom-go/internal/inventory"
	"github.com/sushanthach12/ecom-go/internal/products"
)

func Register(router *chi.Mux, db *database.DB, logger *slog.Logger, productSvc products.Adapter, inventorySvc inventory.Adapter) {
	repository := newRepository(db)
	service := NewService(repository, productSvc, inventorySvc)
	handler := newHandler(logger, service)

	router.Route("/orders", func(r chi.Router) {
		r.Use(auth.RequireAuth)

		r.Get("/", handler.GetOrders)
		r.Post("/", handler.PlaceOrder)
	})
}
