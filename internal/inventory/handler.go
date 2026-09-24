package inventory

import (
	"log/slog"
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/httpx"
	"github.com/sushanthach12/ecom-go/internal/json"
)

type handler struct {
	service Service
	logger  *slog.Logger
}

func newHandler(svc Service, logger *slog.Logger) *handler {
	return &handler{
		service: svc,
		logger:  logger,
	}
}

func (h *handler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info("Received Request for inventory update")

	payload, ok := json.DecodeAndValidate[updateStockDto](w, r, h.logger)
	if !ok {
		return
	}

	result, err := h.service.UpdateStock(ctx, updateStockParams{
		ProductId: payload.ProductId,
		Quantity:  payload.Quantity,
	})
	if err != nil {
		h.logger.Error("Failed to update stock:", "error", err)
		httpx.HandleError(w, err)
		return
	}

	response := constants.NewResponse(result, nil)
	httpx.WriteJSON(w, http.StatusAccepted, response)
}
