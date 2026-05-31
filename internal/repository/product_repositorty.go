package repository

import (
	"context"

	"main.go/internal/domain"
	"main.go/internal/dto"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductByID(ctx context.Context, productID string) (*domain.Product, error)
	Update(ctx context.Context, productID string, req *dto.RequestPatchProduct) error
}
