package inventory

import (
	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/helpers"
)

type updateStockDto struct {
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

func (dto *updateStockDto) Validate() error {
	if helpers.CheckIfStringEmpty(dto.ProductId) {
		return &constants.ValidationError{
			Field:   "productId",
			Message: "must not be empty",
		}
	}

	if !helpers.CheckIfUuid(dto.ProductId) {
		return &constants.ValidationError{
			Field:   "productId",
			Message: "must be a valid id",
		}
	}

	if dto.Quantity <= 0 {
		return &constants.ValidationError{
			Field:   "quantity",
			Message: "must greater than 0",
		}
	}

	return nil
}
