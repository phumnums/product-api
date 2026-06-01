package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"main.go/internal/domain"
	"main.go/internal/dto"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := "postgres://root:pAssw0rd@localhost:5432/product_db?sslmode=disable"

	db, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)

	err = db.Ping(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})
	return db
}

func TestProductPostgresRepository_Create_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewProductPostgresRepository(db)

	description := "Top-To-Toe Hair&Body Bath 100 ml. x48"
	salePrice := float64(1278)

	product := &domain.Product{
		ID:          ulid.Make().String(),
		Name:        "Johnson's Baby",
		Description: &description,
		SalePrice:   &salePrice,
		Price:       1680,
	}

	created, err := repo.Create(context.Background(), product)

	require.NoError(t, err)
	require.NotNil(t, product)

	assert.Equal(t, product.ID, created.ID)
	assert.Equal(t, "Johnson's Baby", created.Name)
	assert.Equal(t, "Top-To-Toe Hair&Body Bath 100 ml. x48", *created.Description)
	assert.Equal(t, float64(1278), *created.SalePrice)
	assert.Equal(t, float64(1680), created.Price)
}

func TestProductPostgresRepository_GetProductByID_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewProductPostgresRepository(db)

	product := &domain.Product{
		ID:    ulid.Make().String(),
		Name:  "Johnson's Baby",
		Price: 1680,
	}

	created, err := repo.Create(context.Background(), product)
	require.NoError(t, err)

	found, err := repo.GetProductByID(context.Background(), created.ID)
	require.NotNil(t, found)
	require.NoError(t, err)

	assert.Equal(t, created.ID, product.ID)
	assert.Equal(t, "Johnson's Baby", created.Name)
	assert.Equal(t, float64(1680), created.Price)
}

func TestProductPostgresRepository_GetProductByID_NotFound(t *testing.T) {
	db := setupTestDB(t)

	repo := NewProductPostgresRepository(db)

	found, err := repo.GetProductByID(context.Background(), "DATANOTFOUND")
	require.Nil(t, found)
	require.Error(t, err)

	assert.Equal(t, err, pgx.ErrNoRows)
}

func TestProductPostgresRepository_Update_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductPostgresRepository(db)

	description := "Top-To-Toe Hair&Body Bath 100 ml. x48"
	salePrice := float64(1278)

	product := &domain.Product{
		ID:          ulid.Make().String(),
		Name:        "Johnson's Baby",
		Description: &description,
		SalePrice:   &salePrice,
		Price:       1680,
	}

	created, err := repo.Create(context.Background(), product)
	require.NoError(t, err)

	newName := "New Name"
	newDescription := "New Description"
	newSalePrice := float64(1200)
	newPrice := float64(1700)
	req := dto.RequestPatchProduct{
		Name: dto.OptionalString{
			Set:   true,
			Value: &newName,
		},
		Description: dto.OptionalString{
			Set:   true,
			Value: &newDescription,
		},
		SalePrice: dto.OptionalFloat{
			Set:   true,
			Value: &newSalePrice,
		},
		Price: dto.OptionalFloat{
			Set:   true,
			Value: &newPrice,
		},
	}

	err = repo.Update(context.Background(), created.ID, &req)
	require.NoError(t, err)

	updated, err := repo.GetProductByID(context.Background(), created.ID)
	require.NoError(t, err)

	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "New Description", *updated.Description)
	assert.Equal(t, float64(1200), *updated.SalePrice)
	assert.Equal(t, float64(1700), updated.Price)
}

func TestProductPostgresRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductPostgresRepository(db)

	newName := "New Name"
	req := dto.RequestPatchProduct{
		Name: dto.OptionalString{
			Set:   true,
			Value: &newName,
		},
	}

	err := repo.Update(context.Background(), "TESTPATCHPRODUCT", &req)
	require.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}
