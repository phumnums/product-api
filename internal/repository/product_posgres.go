package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"main.go/internal/domain"
)

type productPostgresRepository struct {
	db *pgxpool.Pool
}

func NewProductPostgresRepository(db *pgxpool.Pool) ProductRepository {
	return &productPostgresRepository{
		db: db,
	}
}

func (r *productPostgresRepository) Create(ctx context.Context, product domain.Product) (*domain.Product, error) {
	var p domain.Product

	query := `
		INSERT INTO products (
			id, name, description, sale_price, price, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, name, description, sale_price, price, created_at, updated_at
	`
	if err := r.db.QueryRow(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.SalePrice,
		product.Price,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.SalePrice,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &p, nil
}
