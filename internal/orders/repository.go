package orders

import (
	"context"

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

func (r *repository) Create(ctx context.Context, data placeOrderRepoParams) (order, error) {
	exec := database.GetExecutor(ctx, r.db)

	var o order
	row := exec.QueryRowContext(ctx,
		`INSERT INTO orders (customer_id, total)
		 VALUES ($1, $2)
		 RETURNING id, customer_id, status, total, created_at, updated_at`,
		data.CustomerId, data.Total,
	)

	if err := row.Scan(&o.ID, &o.CustomerId, &o.Status, &o.Total, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return order{}, err
	}

	o.Items = make([]orderItems, 0, len(data.Items))

	for _, item := range data.Items {
		var oi orderItems
		itemRow := exec.QueryRowContext(ctx,
			`INSERT INTO order_items (order_id, product_id, price, quantity)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id, order_id, product_id, price, quantity`,
			o.ID, item.ProductId, item.Price, item.Quantity,
		)

		if err := itemRow.Scan(&oi.ID, &oi.OrderId, &oi.ProductId, &oi.Price, &oi.Quantity); err != nil {
			return order{}, err
		}

		o.Items = append(o.Items, oi)
	}

	return o, nil
}
