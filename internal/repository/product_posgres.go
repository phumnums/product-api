package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"main.go/internal/domain"
	"main.go/internal/dto"
)

type productPostgresRepository struct {
	db *pgxpool.Pool
}

func NewProductPostgresRepository(db *pgxpool.Pool) ProductRepository {
	return &productPostgresRepository{
		db: db,
	}
}

func (r *productPostgresRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
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

func (r *productPostgresRepository) GetProductByID(ctx context.Context, productID string) (*domain.Product, error) {
	var p domain.Product

	query := `
		SELECT 
			id, name, description, sale_price, price, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	if err := r.db.QueryRow(ctx, query, productID).Scan(
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

func (r *productPostgresRepository) Update(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
	query := `
		UPDATE products
		SET
			name = CASE WHEN $1 THEN $2 ELSE name END,
			description = CASE WHEN $3 THEN $4 ELSE description END,
			sale_price = CASE WHEN $5 THEN $6 ELSE sale_price END,
			price = CASE WHEN $7 THEN $8 ELSE price END,
			updated_at = NOW()
		WHERE id = $9
	`
	_, err := r.db.Exec(
		ctx,
		query,
		req.Name.Set,
		req.Name.Value,
		req.Description.Set,
		req.Description.Value,
		req.SalePrice.Set,
		req.SalePrice.Value,
		req.Price.Set,
		req.Price.Value,
		productID,
	)
	if err != nil {
		return err
	}

	return nil
}
