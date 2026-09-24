package orders

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/sushanthach12/ecom-go/internal/database"
)

type repository struct {
	db *database.DB
}

func newRepository(db *database.DB) orderRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.TxManager.WithTx(ctx, fn)
}

func (r *repository) Create(ctx context.Context, data placeOrderRepoParams) (orderEntity, error) {
	exec := database.GetExecutor(ctx, r.db)

	var o orderEntity
	row := exec.QueryRowContext(ctx,
		`INSERT INTO orders (customer_id, total)
		 VALUES ($1, $2)
		 RETURNING id, customer_id, status, total, created_at, updated_at`,
		data.CustomerId, data.Total,
	)

	if err := row.Scan(&o.ID, &o.CustomerId, &o.Status, &o.Total, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return orderEntity{}, err
	}

	o.Items = make([]orderItemsEntity, 0, len(data.Items))

	for _, item := range data.Items {
		var oi orderItemsEntity
		itemRow := exec.QueryRowContext(ctx,
			`INSERT INTO order_items (order_id, product_id, product_name, price, quantity)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id, order_id, product_id, product_name, price, quantity`,
			o.ID, item.ProductId, item.ProductName, item.Price, item.Quantity,
		)

		if err := itemRow.Scan(&oi.ID, &oi.OrderId, &oi.ProductId, &oi.ProductName, &oi.Price, &oi.Quantity); err != nil {
			return orderEntity{}, err
		}

		o.Items = append(o.Items, oi)
	}

	return o, nil
}

func (r *repository) Count(ctx context.Context) (int32, error) {
	exec := database.GetExecutor(ctx, r.db)

	var totalItems int32
	if err := exec.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders`).Scan(&totalItems); err != nil {
		return 0, fmt.Errorf("count orders: %w", err)
	}

	return totalItems, nil
}

func (r *repository) GetOrders(ctx context.Context, params listOrdersRepoParams) ([]orderEntity, error) {
	exec := database.GetExecutor(ctx, r.db)

	offset := (params.Page - 1) * params.PageSize

	// 1. Page the order IDs first
	rows, err := exec.QueryContext(ctx, `
		SELECT id, customer_id, total, status, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, params.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	ordersById := make(map[string]*orderEntity)
	var orderIds []string

	for rows.Next() {
		var o orderEntity
		if err := rows.Scan(&o.ID, &o.CustomerId, &o.Total, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan order row: %w", err)
		}
		o.Items = []orderItemsEntity{}
		ordersById[o.ID] = &o
		orderIds = append(orderIds, o.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order rows: %w", err)
	}

	if len(orderIds) == 0 {
		return []orderEntity{}, nil
	}

	// 2. Fetch items only for those order IDs — no LIMIT/OFFSET needed here
	itemRows, err := exec.QueryContext(ctx, `
		SELECT oi.id, oi.order_id, oi.product_id, p.name, oi.quantity, oi.price
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = ANY($1)
	`, pq.Array(orderIds))
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var oi orderItemsEntity
		if err := itemRows.Scan(&oi.ID, &oi.OrderId, &oi.ProductId, &oi.ProductName, &oi.Quantity, &oi.Price); err != nil {
			return nil, fmt.Errorf("scan order item row: %w", err)
		}
		if o, ok := ordersById[oi.OrderId]; ok {
			o.Items = append(o.Items, oi)
		}
	}

	if err := itemRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order item rows: %w", err)
	}

	// preserve original page order
	result := make([]orderEntity, 0, len(orderIds))
	for _, id := range orderIds {
		result = append(result, *ordersById[id])
	}

	return result, nil
}
