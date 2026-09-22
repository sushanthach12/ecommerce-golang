package products

import (
	"log/slog"
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/httpx"
)

type handler struct {
	service Service
	logger  *slog.Logger
}

// constructor
func newHandler(svc Service, logger *slog.Logger) *handler {
	return &handler{
		service: svc,
		logger:  logger,
	}
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Received request for product listing")

	response, err := h.service.List(r.Context())
	if err != nil {
		h.logger.Error("list listings failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, response)
}
