package products

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/sushanthach12/ecom-go/internal/database"
)

type repository struct {
	db *database.DB
}

func newRepository(db *database.DB) productRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Count(ctx context.Context) (int32, error) {
	var totalItems int32
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products`).Scan(&totalItems)
	if err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}
	return totalItems, nil
}

func (r *repository) List(ctx context.Context, limit, skip int32) ([]productEntity, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, price, quantity, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	products := []productEntity{}
	for rows.Next() {
		var l productEntity
		if err := rows.Scan(&l.ID, &l.Name, &l.Price, &l.Quantity, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan product row: %w", err)
		}
		products = append(products, l)
	}

	// rows.Err() must be checked AFTER the loop, not before — it only
	// reflects errors encountered during iteration, not from QueryContext
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product rows: %w", err)
	}

	return products, nil
}

func (r *repository) GetById(ctx context.Context, id string) (productEntity, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, price, quantity, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id)
	if row.Err() != nil {
		return productEntity{}, fmt.Errorf("query products: %w", row.Err())
	}

	var product productEntity
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.CreatedAt, &product.UpdatedAt); err != nil {
		return productEntity{}, fmt.Errorf("scan product row: %w", err)
	}

	return product, nil
}

func (r *repository) FindByIds(ctx context.Context, ids []string) ([]productEntity, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, price, quantity, created_at, updated_at
		FROM products
		WHERE id = ANY($1)
	`, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	var result []productEntity

	for rows.Next() {
		var product productEntity
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product row: %w", err)
		}
		result = append(result, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product rows: %w", err)
	}

	return result, nil
}
