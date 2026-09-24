package inventory

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrInventoryNotFound = errors.New("inventory not found")
)

type adapter struct {
	repo   inventoryRepository
	logger *slog.Logger
}

func newAdapter(repo inventoryRepository, logger *slog.Logger) Adapter {
	return &adapter{
		repo:   repo,
		logger: logger,
	}
}

func (a *adapter) GetByProductID(ctx context.Context, id string) (Inventory, error) {
	inventory, err := a.repo.GetByProductId(ctx, id)
	if err != nil {
		a.logger.Error("adapter: GetByProductID", "error", err)
		return Inventory{}, ErrInventoryNotFound
	}

	return Inventory{
		ID:        inventory.ID,
		ProductId: inventory.ProductId,
		Quantity:  inventory.Quantity,
		CreatedAt: inventory.CreatedAt,
		UpdatedAt: inventory.UpdatedAt,
	}, nil
}

/*
 * Products needs to be checked before calling this method
 */
func (s *adapter) DecrementStock(ctx context.Context, payload DecrementStockParam) error {
	if payload.ProductId == "" {
		return fmt.Errorf("productId is required")
	}

	if payload.Quantity <= 0 {
		return fmt.Errorf("quantity is required")
	}

	if err := s.repo.DecrementStock(ctx, stockDecrementRepoParam{
		ProductId: payload.ProductId,
		Quantity:  payload.Quantity,
	}); err != nil {
		return fmt.Errorf("service: decrement stock: %w", err)
	}

	return nil
}

/*
 * Products needs to be checked before calling this method
 */
func (a *adapter) DecrementStockBulk(ctx context.Context, payload []DecrementStockParam) error {
	if len(payload) == 0 {
		return fmt.Errorf("payload is required")
	}

	repoParams := make([]stockDecrementRepoParam, 0, len(payload))
	seen := make(map[string]bool, len(payload))

	for i, item := range payload {
		if item.ProductId == "" {
			return fmt.Errorf("item %d: productId is required", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity is required", i)
		}
		if seen[item.ProductId] {
			return fmt.Errorf("item %d: duplicate productId %s in payload", i, item.ProductId)
		}
		seen[item.ProductId] = true

		repoParams = append(repoParams, stockDecrementRepoParam{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	if err := a.repo.DecrementStockBulk(ctx, repoParams); err != nil {
		return fmt.Errorf("service: decrement stock bulk: %w", err)
	}

	return nil
}
