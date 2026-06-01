package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"main.go/internal/domain"
	"main.go/internal/dto"
)

type mockProductRepository struct {
	CreateFunc         func(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductByIDFunc func(ctx context.Context, productID string) (*domain.Product, error)
	UpdateFunc         func(ctx context.Context, productID string, req *dto.RequestPatchProduct) error
}

func (m *mockProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	return m.CreateFunc(ctx, product)
}

func (m *mockProductRepository) GetProductByID(ctx context.Context, productID string) (*domain.Product, error) {
	return m.GetProductByIDFunc(ctx, productID)
}

func (m *mockProductRepository) Update(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
	return m.UpdateFunc(ctx, productID, req)
}

func TestCreateProductService_Success(t *testing.T) {
	repo := &mockProductRepository{
		CreateFunc: func(ctx context.Context, product *domain.Product) (*domain.Product, error) {
			return product, nil
		},
	}

	svc := NewProductService(repo)

	description := "Top-To-Toe Hair&Body Bath 100 ml. x48"
	salePrice := float64(1278)

	req := &dto.RequestCreateProduct{
		Name:        "Johnson's Baby",
		Description: &description,
		SalePrice:   &salePrice,
		Price:       1680,
	}

	res, status, err := svc.CreateProductService(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Johnson's Baby", res.Name)
	assert.NotNil(t, "Top-To-Toe Hair&Body Bath 100 ml. x48", *res.Description)
	assert.Equal(t, float64(1278), *res.SalePrice)
	assert.Equal(t, float64(1680), res.Price)
}

func TestCreateProductService_Success_WhenOptionalFieldsAreNil(t *testing.T) {
	repo := &mockProductRepository{
		CreateFunc: func(ctx context.Context, product *domain.Product) (*domain.Product, error) {
			return product, nil
		},
	}

	svc := NewProductService(repo)

	req := &dto.RequestCreateProduct{
		Name:  "Johnson's Baby",
		Price: 1680,
	}

	res, status, err := svc.CreateProductService(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Johnson's Baby", res.Name)
	assert.Nil(t, res.Description)
	assert.Nil(t, res.SalePrice)
	assert.Equal(t, float64(1680), res.Price)
}

func TestCreateProductService_BlankDescription_ShouldBeNil(t *testing.T) {
	repo := &mockProductRepository{
		CreateFunc: func(ctx context.Context, product *domain.Product) (*domain.Product, error) {
			return product, nil
		},
	}

	svc := NewProductService(repo)

	description := "      "

	req := &dto.RequestCreateProduct{
		Name:        "Johnson's Baby",
		Description: &description,
		Price:       1680,
	}

	res, status, err := svc.CreateProductService(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Johnson's Baby", res.Name)
	assert.Nil(t, res.Description)
	assert.Equal(t, float64(1680), res.Price)
}

func TestCreateProductService_SalePriceGreaterThanPrice(t *testing.T) {
	repo := &mockProductRepository{}
	svc := NewProductService(repo)

	salePrice := float64(2000)

	req := &dto.RequestCreateProduct{
		Name:      "Johnson's Baby",
		SalePrice: &salePrice,
		Price:     1680,
	}

	res, status, err := svc.CreateProductService(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, res)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "sale_price must be lower than price", err.Error())
}

func TestCreateProductService_RepositoryError(t *testing.T) {
	repo := &mockProductRepository{
		CreateFunc: func(ctx context.Context, product *domain.Product) (*domain.Product, error) {
			return nil, errors.New("database error")
		},
	}

	svc := NewProductService(repo)

	req := &dto.RequestCreateProduct{
		Name:  "Johnson's Baby",
		Price: 1680,
	}

	res, status, err := svc.CreateProductService(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, res)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, "server error", err.Error())
}
