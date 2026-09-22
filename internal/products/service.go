package products

import (
	"context"
	"fmt"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

type service struct {
	repo productRepository
}

func NewService(repo productRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) List(ctx context.Context, params listProductPayload) (constants.Response[listProductsResponseDto], error) {
	page := params.Page
	limit := params.Limit

	skip := (page - 1) * limit

	totalItems, err := s.repo.Count(ctx)
	if err != nil {
		return constants.Response[listProductsResponseDto]{}, fmt.Errorf("service: get product count: %w", err)
	}

	products, err := s.repo.List(ctx, limit, skip)
	if err != nil {
		return constants.Response[listProductsResponseDto]{}, fmt.Errorf("service: list products: %w", err)
	}

	items := make([]listProductsResponseDto, 0, len(products))
	for _, p := range products {
		items = append(items, listProductsResponseDto{
			ID:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			Quantity:  p.Quantity,
			CreatedAt: p.CreatedAt,
		})
	}

	totalPages := (totalItems + limit - 1) / limit

	return constants.NewPaginatedResponse(items, constants.Pagination{
		Page:       page,
		PageSize:   limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}), nil
}

func (s *service) GetById(ctx context.Context, id string) (constants.SimpleResponse[Product], error) {
	product, err := s.repo.GetById(ctx, id)
	if err != nil {
		return constants.NewResponse(Product{}, nil), fmt.Errorf("service: list products: %w", err)
	}

	return constants.NewResponse(Product{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Quantity:  product.Quantity,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil), nil
}
