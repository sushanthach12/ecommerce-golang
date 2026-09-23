package products

import (
	"log/slog"
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/helpers"
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
	params := helpers.ParsePaginationParams(r)

	h.logger.Info("Received request for product listing")

	response, err := h.service.List(r.Context(), listProductPayload{
		Page:  params.Page,
		Limit: params.Limit,
	})
	if err != nil {
		h.logger.Error("list products failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, response)
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	productId := helpers.GetPathValue(r, "product_id")

	h.logger.Info("Received request for product details")

	if productId == "" {
		h.logger.Error("Invalid product id")
		httpx.Error(w, http.StatusBadGateway, "Invalid Product Id", httpx.CodeInvalidId)
		return
	}

	response, err := h.service.GetById(r.Context(), productId)
	if err != nil {
		if err == ErrProductNotFound {
			h.logger.Error("product not found", "error", err)
			httpx.Error(w, http.StatusNotFound, "Product not found!", httpx.CodeNotFound)
			return
		}
		h.logger.Error("product details failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, response)
}
