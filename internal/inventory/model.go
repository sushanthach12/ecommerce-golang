package inventory

import "time"

type inventoryEntity struct {
	ID        string
	ProductId string
	Quantity  int32

	CreatedAt time.Time
	UpdatedAt time.Time
}
