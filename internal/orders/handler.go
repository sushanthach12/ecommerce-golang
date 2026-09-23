package orders

import (
	"log/slog"
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/helpers"
	"github.com/sushanthach12/ecom-go/internal/httpx"
	"github.com/sushanthach12/ecom-go/internal/json"
)

type handler struct {
	service Service
	logger  *slog.Logger
}

func newHandler(logger *slog.Logger, service Service) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	params := helpers.ParsePaginationParams(r)

	h.logger.Info("Received request for orders listing")

	response, err := h.service.List(r.Context(), listOrderParams{
		Page:  params.Page,
		Limit: params.Limit,
	})
	if err != nil {
		h.logger.Error("list orders failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, response)

}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info("Received Request for place order")

	payload, ok := json.DecodeAndValidate[placeOrderPayloadDto](w, r, h.logger)
	if !ok {
		return
	}

	result, err := h.service.PlaceOrder(ctx, placeOrderParams{
		CustomerId: payload.CustomerId,
		Items:      mapItems(payload.Items),
	})
	if err != nil {
		h.logger.Error("Failed to place order:", "error", err)
		httpx.HandleError(w, err)
		return
	}

	response := constants.NewResponse(result, nil)
	httpx.WriteJSON(w, http.StatusCreated, response)
}
