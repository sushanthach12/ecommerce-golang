package products

import (
	"context"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

type Service interface {
	List(ctx context.Context, params listProductPayload) (constants.Response[listProductsResponseDto], error)
}

type ProductRepository interface {
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, skip int) ([]productEntity, error)
}
