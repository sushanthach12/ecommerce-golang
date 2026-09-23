package orders

import (
	"fmt"

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
