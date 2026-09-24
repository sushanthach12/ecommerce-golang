package inventory

import (
	"context"
	"time"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

type Inventory struct {
	ID        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type updateStockParams struct {
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type updateStockResponse struct {
	ProductId string `json:"product_id"`
}

type DecrementStockParam struct {
	ProductId string
	Quantity  int32
}

type Service interface {
	UpdateStock(ctx context.Context, paylaod updateStockParams) (constants.SimpleResponse[updateStockResponse], error)
}

// ================= Repository==========================
type stockDecrementRepoParam struct {
	ProductId string
	Quantity  int32
}

type createRepoParam struct {
	ProductId string
	Quantity  int32
}

type updateRepoParam struct {
	ProductId string
	Quantity  int32
}

type inventoryRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	Create(ctx context.Context, data createRepoParam) (inventoryEntity, error)
	Update(ctx context.Context, data updateRepoParam) error
	GetByProductId(ctx context.Context, id string) (inventoryEntity, error)
	DecrementStock(ctx context.Context, data stockDecrementRepoParam) error
	DecrementStockBulk(ctx context.Context, items []stockDecrementRepoParam) error
}

// =================== Adapter===================

type Adapter interface {
	GetByProductID(ctx context.Context, id string) (Inventory, error)
	DecrementStock(ctx context.Context, payload DecrementStockParam) error
	DecrementStockBulk(ctx context.Context, payload []DecrementStockParam) error
}
