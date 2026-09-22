package products

import "time"

type listProductPayload struct {
	Page  int
	Limit int
}

type listProductsResponseDto struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float32   `json:"price"`
	Quantity  int32     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
