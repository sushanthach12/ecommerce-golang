package orders

import "time"

type orderItemsEntity struct {
	ID          string
	OrderId     string
	ProductId   string
	ProductName string
	Quantity    int32
	Price       float64
}

type orderEntity struct {
	ID         string
	CustomerId string
	Total      float64
	Status     string
	Items      []orderItemsEntity

	CreatedAt time.Time
	UpdatedAt time.Time
}
