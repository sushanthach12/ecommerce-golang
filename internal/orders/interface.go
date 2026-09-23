package orders

import (
	"context"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

type placeOrderItemParam struct {
	ProductId string
	Quantity  int32
}

type placeOrderParams struct {
	CustomerId string
	Items      []placeOrderItemParam
}

func mapItems(items []placeOrderItemsPayloadDto) []placeOrderItemParam {
	result := make([]placeOrderItemParam, len(items))

	for i, item := range items {
		result[i] = placeOrderItemParam{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}

	return result
}

type Service interface {
	List(ctx context.Context) (constants.Response[listOrderResponseDto], error)
	PlaceOrder(ctx context.Context, payload placeOrderParams) (constants.SimpleResponse[placeOrderResponseDto], error)
}

type placeOrderItemRepoParam struct {
	ProductId   string
	ProductName string
	Price       float64
	Quantity    int32
}

type placeOrderRepoParams struct {
	CustomerId string
	Total      float64
	Items      []placeOrderItemRepoParam
}

type orderRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	Create(ctx context.Context, data placeOrderRepoParams) (order, error)
}
