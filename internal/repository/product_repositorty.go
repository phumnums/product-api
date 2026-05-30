package repository

import (
	"context"

	"main.go/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) (*domain.Product, error)
}
