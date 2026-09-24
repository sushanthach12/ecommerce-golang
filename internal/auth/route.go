package auth

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/sushanthach12/ecom-go/internal/users"
)

func Register(router *chi.Mux, logger *slog.Logger, userSvc users.Service) {
	service := newService(userSvc, logger)
	handler := newHandler(service, logger)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}
