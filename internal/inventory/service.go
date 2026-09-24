package inventory

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/httpx"
	"github.com/sushanthach12/ecom-go/internal/products"
)

var (
	ErrNotFound = errors.New("inventory record not found")
)

type service struct {
	logger *slog.Logger

	repo           inventoryRepository
	productAdapter products.Adapter
}

func NewService(logger *slog.Logger, repo inventoryRepository, productAdapter products.Adapter) Service {
	return &service{
		logger:         logger,
		repo:           repo,
		productAdapter: productAdapter,
	}
}

func (s *service) UpdateStock(ctx context.Context, payload updateStockParams) (constants.SimpleResponse[updateStockResponse], error) {
	var defaultResponse constants.SimpleResponse[updateStockResponse]

	// check if the product exists
	productDetails, err := s.productAdapter.GetByID(ctx, payload.ProductId)
	if err != nil {
		return defaultResponse, err
	}

	if productDetails.ID == "" {
		return defaultResponse, fmt.Errorf("invalid product")
	}

	if err := s.repo.Update(ctx, updateRepoParam{
		ProductId: payload.ProductId,
		Quantity:  payload.Quantity,
	}); err != nil {
		if errors.Is(err, ErrNotFound) {
			return defaultResponse, &httpx.NotFoundError{
				Message: "please create the inventory for the product",
			}
		}
		return defaultResponse, fmt.Errorf("Something went wrong!")
	}

	return constants.NewResponse(updateStockResponse{}, nil), nil
}
