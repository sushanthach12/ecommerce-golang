package health

import (
	"github.com/go-chi/chi/v5"
)

func Register(router *chi.Mux) {
	handler := newHandler()

	router.Get("/health", handler.Health)
}
