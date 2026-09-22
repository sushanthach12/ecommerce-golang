package products

import (
	"context"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) List(ctx context.Context) (constants.Response[getProductsResponseDto], error) {

	products := make([]getProductsResponseDto, 0)

	products = append(products, getProductsResponseDto{
		ID:   "product-1",
		Name: "product",
	})

	return constants.NewPaginatedResponse(products, constants.Pagination{
		Page:       1,
		PageSize:   10,
		TotalItems: 10,
		TotalPages: 1,
	}), nil
}
