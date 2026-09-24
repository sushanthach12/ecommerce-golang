package users

import (
	"log/slog"
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/httpx"
)

type handler struct {
	service Service
	logger  *slog.Logger
}

func newHandler(service Service, logger *slog.Logger) handler {
	return handler{service: service, logger: logger}
}

func (h *handler) Me(w http.ResponseWriter, r *http.Request) {

	httpx.WriteJSON(w, http.StatusCreated, "me")
}
