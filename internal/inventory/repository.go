package inventory

import (
	"context"
	"fmt"

	"github.com/sushanthach12/ecom-go/internal/database"
)

type repository struct {
	db *database.DB
}

func newRepository(db *database.DB) inventoryRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.TxManager.WithTx(ctx, fn)
}

func (r *repository) Create(ctx context.Context, data createRepoParam) (inventoryEntity, error) {
	exec := database.GetExecutor(ctx, r.db)

	var inventory inventoryEntity
	row := exec.QueryRowContext(ctx,
		`INSERT INTO inventory (product_id, quantity)
		 VALUES ($1, $2)
		 RETURNING id, product_id, quantity, created_at, updated_at`,
		data.ProductId, data.Quantity,
	)

	if err := row.Scan(&inventory.ID, &inventory.ProductId, &inventory.Quantity, &inventory.CreatedAt, &inventory.UpdatedAt); err != nil {
		return inventoryEntity{}, err
	}

	return inventory, nil
}

func (r *repository) Update(ctx context.Context, data updateRepoParam) error {
	exec := database.GetExecutor(ctx, r.db)

	result, err := exec.ExecContext(ctx, `
			UPDATE inventory
			SET quantity = $1, updated_at = NOW()
			WHERE product_id = $2
		`, data.Quantity, data.ProductId)
	if err != nil {
		return fmt.Errorf("update stock for product %s: %w", data.ProductId, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected for inventory %s: %w", data.ProductId, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product %s: %w", data.ProductId, ErrNotFound)
	}

	return nil
}

func (r *repository) GetByProductId(ctx context.Context, id string) (inventoryEntity, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, product_id, quantity, created_at, updated_at
		FROM inventory
		WHERE product_id = $1
	`, id)
	if row.Err() != nil {
		return inventoryEntity{}, fmt.Errorf("query inventory: %w", row.Err())
	}

	var inventory inventoryEntity
	if err := row.Scan(&inventory.ID, &inventory.ProductId, &inventory.Quantity, &inventory.CreatedAt, &inventory.UpdatedAt); err != nil {
		return inventoryEntity{}, fmt.Errorf("scan inventory row: %w", err)
	}

	return inventory, nil
}

func (r *repository) DecrementStock(ctx context.Context, data stockDecrementRepoParam) error {
	exec := database.GetExecutor(ctx, r.db)

	result, err := exec.ExecContext(ctx, `
			UPDATE inventory
			SET quantity = quantity - $1, updated_at = NOW()
			WHERE product_id = $2 AND quantity >= $1
		`, data.Quantity, data.ProductId)
	if err != nil {
		return fmt.Errorf("decrement stock for product %s: %w", data.ProductId, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected for inventory %s: %w", data.ProductId, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient stock or product %s not found", data.ProductId)
	}

	return nil
}

func (r *repository) DecrementStockBulk(ctx context.Context, items []stockDecrementRepoParam) error {
	exec := database.GetExecutor(ctx, r.db)

	for _, item := range items {
		result, err := exec.ExecContext(ctx, `
			UPDATE inventory
			SET quantity = quantity - $1, updated_at = NOW()
			WHERE product_id = $2 AND quantity >= $1
		`, item.Quantity, item.ProductId)
		if err != nil {
			return fmt.Errorf("decrement stock for product %s: %w", item.ProductId, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("check rows affected for product %s: %w", item.ProductId, err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("insufficient stock or product %s not found", item.ProductId)
		}
	}

	return nil
}
