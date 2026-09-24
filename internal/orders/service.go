package orders

import (
	"context"
	"fmt"
	"sync"

	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/httpx"
	"github.com/sushanthach12/ecom-go/internal/inventory"
	"github.com/sushanthach12/ecom-go/internal/products"
)

type service struct {
	repo         orderRepository
	productSvc   products.Adapter
	inventorySvc inventory.Adapter
}

func NewService(repo orderRepository, productSvc products.Adapter, inventorySvc inventory.Adapter) Service {
	return &service{
		repo:         repo,
		productSvc:   productSvc,
		inventorySvc: inventorySvc,
	}
}

func (s *service) PlaceOrder(ctx context.Context, payload placeOrderParams) (constants.SimpleResponse[placeOrderResponseDto], error) {
	defaultResponse := constants.SimpleResponse[placeOrderResponseDto]{}

	// 1. Look up all products referenced in the order in one batch call
	productIds := make([]string, len(payload.Items))
	for i, item := range payload.Items {
		productIds[i] = item.ProductId
	}

	productsData, err := s.productSvc.FindByIds(ctx, productIds)
	if err != nil {
		return defaultResponse, fmt.Errorf("failed to fetch products: %w", err)
	}

	productById := make(map[string]products.Product, len(productsData))
	for _, p := range productsData {
		productById[p.ID] = p
	}

	// 2. Validate each item: product exists, has enough stock
	orderItems := make([]placeOrderItemRepoParam, 0, len(payload.Items))
	var total float64

	for _, item := range payload.Items {
		product, ok := productById[item.ProductId]
		if !ok {
			return defaultResponse, &httpx.NotFoundError{
				Message: fmt.Sprintf("product %s not found", item.ProductId),
			}
		}

		if product.Quantity < item.Quantity {
			return defaultResponse, &httpx.ConflictError{
				Field:   "items.quantity",
				Message: fmt.Sprintf("insufficient stock for product %s", item.ProductId),
			}
		}

		orderItems = append(orderItems, placeOrderItemRepoParam{
			ProductId:   item.ProductId,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
		})

		total += product.Price * float64(item.Quantity)
	}

	var decrementOrderItems []inventory.DecrementStockParam
	for _, item := range orderItems {
		decrementOrderItems = append(decrementOrderItems, inventory.DecrementStockParam{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	// 3 & 4. Create the order and decrement stock atomically
	var createdOrder orderEntity
	err = s.repo.WithTx(ctx, func(txCtx context.Context) error {
		var txErr error
		createdOrder, txErr = s.repo.Create(txCtx, placeOrderRepoParams{
			CustomerId: payload.CustomerId,
			Total:      total,
			Items:      orderItems,
		})
		if txErr != nil {
			return txErr
		}

		return s.inventorySvc.DecrementStockBulk(txCtx, decrementOrderItems)
	})
	if err != nil {
		return defaultResponse, fmt.Errorf("failed to place order: %w", err)
	}

	// 5. Map to response DTO
	return constants.NewResponse(placeOrderResponseDto{ID: createdOrder.ID}, nil), nil
}

func (s *service) List(ctx context.Context, params listOrderParams) (constants.Response[listOrderResponseDto], error) {
	var (
		orders     []orderEntity
		totalItems int32
		ordersErr  error
		countErr   error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		orders, ordersErr = s.repo.GetOrders(ctx, listOrdersRepoParams{
			Page:     params.Page,
			PageSize: params.Limit,
		})
	}()

	go func() {
		defer wg.Done()
		totalItems, countErr = s.repo.Count(ctx)
	}()

	wg.Wait()

	if ordersErr != nil {
		return constants.Response[listOrderResponseDto]{}, fmt.Errorf("service: list orders: %w", ordersErr)
	}
	if countErr != nil {
		return constants.Response[listOrderResponseDto]{}, fmt.Errorf("service: count orders: %w", countErr)
	}

	result := make([]listOrderResponseDto, len(orders))
	for i, o := range orders {
		result[i] = mapOrderToListResponse(o)
	}

	totalPages := (totalItems + params.Limit - 1) / params.Limit

	return constants.NewPaginatedResponse(result, constants.Pagination{
		Page:       params.Page,
		PageSize:   params.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}), nil
}

func mapOrderToListResponse(o orderEntity) listOrderResponseDto {
	items := make([]listOrderItemResponseDto, len(o.Items))
	for i, item := range o.Items {
		items[i] = listOrderItemResponseDto{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		}
	}

	return listOrderResponseDto{
		ID:         o.ID,
		CustomerId: o.CustomerId,
		Total:      o.Total,
		Status:     o.Status,
		Items:      items,
		OrderedAt:  o.CreatedAt,
	}
}
