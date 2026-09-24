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

type listOrderParams struct {
	Page  int32
	Limit int32
}

type Service interface {
	List(ctx context.Context, params listOrderParams) (constants.Response[listOrderResponseDto], error)
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

type listOrdersRepoParams struct {
	Page     int32
	PageSize int32
}

type orderRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	Create(ctx context.Context, data placeOrderRepoParams) (orderEntity, error)
	Count(ctx context.Context) (int32, error)
	GetOrders(ctx context.Context, params listOrdersRepoParams) ([]orderEntity, error)
}
