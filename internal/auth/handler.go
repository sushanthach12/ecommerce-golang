package auth

import (
	"log/slog"
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/httpx"
	"github.com/sushanthach12/ecom-go/internal/json"
)

type handler struct {
	service Service
	logger  *slog.Logger
}

func newHandler(service Service, logger *slog.Logger) handler {
	return handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("received register request")

	payload, ok := json.DecodeAndValidate[registerPayloadDto](w, r, h.logger)
	if !ok {
		return
	}

	response, err := h.service.Register(r.Context(), registerParams{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		h.logger.Error("Failed to register:", "error", err)
		httpx.HandleError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, response)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("received login request")

	payload, ok := json.DecodeAndValidate[loginPayloadDto](w, r, h.logger)
	if !ok {
		return
	}

	response, err := h.service.Login(r.Context(), loginParams{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		h.logger.Error("Failed to login:", "error", err)
		httpx.HandleError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, response)
}
