package products

import (
	"context"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

type Service interface {
	List(ctx context.Context) (constants.Response[getProductsResponseDto], error)
}
