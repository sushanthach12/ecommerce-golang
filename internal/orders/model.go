package orders

import "time"

type orderItems struct {
	ID        string
	OrderId   string
	ProductId string

	// product name

	Quantity int32
	Price    float32
}

type order struct {
	ID         string
	CustomerId string
	Total      float64
	Status     string
	Items      []orderItems

	CreatedAt time.Time
	UpdatedAt time.Time
}
