package products

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
)

func Register(router *chi.Mux, db *sql.DB, logger *slog.Logger) {
	service := NewService()
	handler := newHandler(service, logger)

	router.Get("/products", handler.List)
}
