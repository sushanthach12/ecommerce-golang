package products

import (
	"context"
	"fmt"
	"log/slog"
)

type adapter struct {
	repo   productRepository
	logger *slog.Logger
}

func newAdapter(repo productRepository, logger *slog.Logger) Adapter {
	return &adapter{
		repo:   repo,
		logger: logger,
	}
}

func (a *adapter) GetByID(ctx context.Context, id string) (Product, error) {
	product, err := a.repo.GetById(ctx, id)
	if err != nil {
		return Product{}, fmt.Errorf("service: list products: %w", err)
	}

	return Product{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Quantity:  product.Quantity,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

func (a *adapter) FindByIds(ctx context.Context, productIds []string) ([]Product, error) {
	products, err := a.repo.FindByIds(ctx, productIds)
	if err != nil {
		a.logger.Error("service: list products", "error", err)
		return nil, fmt.Errorf("service: list products: %w", err)
	}

	result := make([]Product, len(products))
	for i, p := range products {
		result[i] = Product{
			ID:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			Quantity:  p.Quantity,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	return result, nil
}

func (a *adapter) DecrementStock(ctx context.Context, items []StockDecrementParam) error {
	// validate items
	var validatedItems []stockDecrementRepoParam

	for i, item := range items {
		if item.ProductId == "" {
			return fmt.Errorf("items[%d].productId is required", i)
		}

		if item.Quantity <= 0 {
			return fmt.Errorf("items[%d].quantity is required", i)
		}

		validatedItems = append(validatedItems, stockDecrementRepoParam{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	if err := a.repo.DecrementStock(ctx, validatedItems); err != nil {
		return fmt.Errorf("service: decrement stock: %w", err)
	}

	return nil
}
