package products

import (
	"context"
	"time"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

// Global interface to used for the product by other modules
type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float32   `json:"price"`
	Quantity  int32     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Service interface {
	List(ctx context.Context, params listProductPayload) (constants.Response[listProductsResponseDto], error)
	GetById(ctx context.Context, id string) (constants.SimpleResponse[Product], error)
}

type productRepository interface {
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, skip int) ([]productEntity, error)
	GetById(ctx context.Context, id string) (productEntity, error)
}

type Adapter interface {
	GetByID(ctx context.Context, id string) (Product, error)
}
