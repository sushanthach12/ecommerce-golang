package orders

import (
	"fmt"
	"time"

	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/helpers"
)

type placeOrderItemsPayloadDto struct {
	ProductId string `json:"productId"`
	Quantity  int32  `json:"quantity"`
}

type placeOrderPayloadDto struct {
	CustomerId string                      `json:"customerId"`
	Items      []placeOrderItemsPayloadDto `json:"items"`
}

func (payload *placeOrderPayloadDto) Validate() error {
	if helpers.CheckIfStringEmpty(payload.CustomerId) {
		return &constants.ValidationError{
			Field:   "customerId",
			Message: "must not be empty",
		}
	}

	if helpers.CheckArrayEmpty(payload.Items) {
		return &constants.ValidationError{
			Field:   "items",
			Message: "must not be empty",
		}
	}

	for i, item := range payload.Items {
		if helpers.CheckIfStringEmpty(item.ProductId) {
			return &constants.ValidationError{
				Field:   fmt.Sprintf("items[%d].productId", i),
				Message: "must not be empty",
			}
		}

		if item.Quantity <= 0 {
			return &constants.ValidationError{
				Field:   fmt.Sprintf("items[%d].quantity", i),
				Message: "must be greater than 0",
			}
		}
	}

	return nil
}

type placeOrderResponseDto struct {
	ID string `json:"id"`
}

type listOrderItemResponseDto struct {
	ProductId   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int32   `json:"quantity"`
}

type listOrderResponseDto struct {
	ID         string                     `json:"id"`
	CustomerId string                     `json:"customer_d"`
	Total      float64                    `json:"total"`
	Status     string                     `json:"status"`
	Items      []listOrderItemResponseDto `json:"items"`
	OrderedAt  time.Time                  `json:"ordered_at"`
}
