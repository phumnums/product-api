package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
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

func TestCreateProductService_BlankName_ShouldReturnBasRequest(t *testing.T) {
	repo := &mockProductRepository{
		CreateFunc: func(ctx context.Context, product *domain.Product) (*domain.Product, error) {
			return product, nil
		},
	}

	svc := NewProductService(repo)

	req := &dto.RequestCreateProduct{
		Name:  "  ",
		Price: 1680,
	}

	res, status, err := svc.CreateProductService(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, res)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "name is required", err.Error())
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

func TestPatchProductService_Success(t *testing.T) {

	id := "TESTPATCHPRODUCT"
	description := "Top-To-Toe Hair&Body Bath 100 ml. x48"
	salePrice := float64(1278)

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)

			return &domain.Product{
				ID:          id,
				Name:        "Johnson's Baby",
				Description: &description,
				SalePrice:   &salePrice,
				Price:       float64(1680),
			}, nil
		},

		UpdateFunc: func(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
			require.NotNil(t, req.Name.Value)
			require.NotNil(t, req.Price.Value)

			assert.Equal(t, id, productID)

			assert.True(t, req.Name.Set)
			assert.Equal(t, "New Name", *req.Name.Value)

			assert.True(t, req.Description.Set)
			assert.Equal(t, "New Description", *req.Description.Value)

			assert.True(t, req.SalePrice.Set)
			assert.Equal(t, float64(1200), *req.SalePrice.Value)

			assert.True(t, req.Price.Set)
			assert.Equal(t, float64(1700), *req.Price.Value)

			return nil
		},
	}

	newName := "New Name"
	newDescription := "New Description"
	newSalePrice := float64(1200)
	newPrice := float64(1700)

	svc := NewProductService(repo)

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

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
}

func TestPatchProductService_Success_NullDescriptionAndSalePrice_ShoudBeNull(t *testing.T) {

	id := "TESTPATCHPRODUCT"
	description := "Top-To-Toe Hair&Body Bath 100 ml. x48"
	salePrice := float64(1278)

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)
			return &domain.Product{
				ID:          id,
				Name:        "Johnson's Baby",
				Description: &description,
				SalePrice:   &salePrice,
				Price:       float64(1680),
			}, nil
		},

		UpdateFunc: func(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
			assert.Equal(t, id, productID)

			assert.True(t, req.Description.Set)
			assert.Nil(t, req.Description.Value)

			assert.True(t, req.SalePrice.Set)
			assert.Nil(t, req.SalePrice.Value)

			return nil
		},
	}

	svc := NewProductService(repo)

	req := dto.RequestPatchProduct{
		Description: dto.OptionalString{
			Set:   true,
			Value: nil,
		},
		SalePrice: dto.OptionalFloat{
			Set:   true,
			Value: nil,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
}

func TestPatchProductService_Success_BlankDescription_ShoudBeNull(t *testing.T) {

	id := "TESTPATCHPRODUCT"
	description := "Top-To-Toe Hair&Body Bath 100 ml. x48"
	salePrice := float64(1278)

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)

			return &domain.Product{
				ID:          id,
				Name:        "Johnson's Baby",
				Description: &description,
				SalePrice:   &salePrice,
				Price:       float64(1680),
			}, nil
		},

		UpdateFunc: func(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
			assert.Equal(t, id, productID)

			assert.True(t, req.Description.Set)
			assert.Nil(t, req.Description.Value)

			return nil
		},
	}

	svc := NewProductService(repo)

	newDescription := "  "

	req := dto.RequestPatchProduct{
		Description: dto.OptionalString{
			Set:   true,
			Value: &newDescription,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
}

func TestPatchProductionService_NullName_ShouldReturnBadRequest(t *testing.T) {
	id := "TESTPATCHPRODUCT"
	isUpdated := false

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)

			return &domain.Product{
				ID:    id,
				Name:  "Johnson's Baby",
				Price: float64(1680),
			}, nil
		},

		UpdateFunc: func(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
			isUpdated = true
			return nil
		},
	}

	svc := NewProductService(repo)

	req := dto.RequestPatchProduct{
		Name: dto.OptionalString{
			Set:   true,
			Value: nil,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "name is required", err.Error())

	assert.False(t, isUpdated)
}

func TestPatchProductionService_BlankName_ShouldReturnBadRequest(t *testing.T) {
	id := "TESTPATCHPRODUCT"
	isUpdated := false

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)
			return &domain.Product{
				ID:    id,
				Name:  "Johnson's Baby",
				Price: float64(1680),
			}, nil
		},

		UpdateFunc: func(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
			isUpdated = true
			return nil
		},
	}

	svc := NewProductService(repo)

	newName := "   "

	req := dto.RequestPatchProduct{
		Name: dto.OptionalString{
			Set:   true,
			Value: &newName,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "name is required", err.Error())

	assert.False(t, isUpdated)
}

func TestPatchProductionService_NullPrice_ShouldReturnBadRequest(t *testing.T) {
	id := "TESTPATCHPRODUCT"
	isUpdated := false

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)

			return &domain.Product{
				ID:    id,
				Name:  "Johnson's Baby",
				Price: float64(1680),
			}, nil
		},

		UpdateFunc: func(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
			isUpdated = true
			return nil
		},
	}

	svc := NewProductService(repo)

	req := dto.RequestPatchProduct{
		Price: dto.OptionalFloat{
			Set:   true,
			Value: nil,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "price is required", err.Error())

	assert.False(t, isUpdated)
}

func TestPatchProductService_EmptyBody_ShouldReturnBadRequest(t *testing.T) {
	id := "TESTPATCHPRODUCT"

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)
			return &domain.Product{
				ID:    id,
				Name:  "Johnson's Baby",
				Price: float64(1680),
			}, nil
		},
	}
	svc := NewProductService(repo)

	req := dto.RequestPatchProduct{}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "invalid request", err.Error())
}

func TestPatchProductService_EmptyProductID_ShouldReturnBadRequest(t *testing.T) {
	repo := &mockProductRepository{}

	svc := NewProductService(repo)

	newName := "New Name"
	req := dto.RequestPatchProduct{
		Name: dto.OptionalString{
			Set:   true,
			Value: &newName,
		},
	}

	status, err := svc.PatchProductService(context.Background(), "", &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "product_id is required", err.Error())
}

func TestPatchProductService_ProductNotFound_ShouldReturnNotFound(t *testing.T) {
	id := "TESTPATCHPRODUCT"

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)
			return nil, pgx.ErrNoRows
		},
	}
	svc := NewProductService(repo)

	newName := "New Name"
	req := dto.RequestPatchProduct{
		Name: dto.OptionalString{
			Set:   true,
			Value: &newName,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "product not found", err.Error())
}

func TestPatchProductService_SalePriceGreaterThanPrice_ShouldReturnBadRequest(t *testing.T) {
	id := "TESTPATCHPRODUCT"
	salePrice := float64(1278)

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)
			return &domain.Product{
				ID:        id,
				Name:      "Johnson's Baby",
				SalePrice: &salePrice,
				Price:     float64(1680),
			}, nil
		},
	}
	svc := NewProductService(repo)

	newSalePrice := float64(1700)
	req := dto.RequestPatchProduct{
		SalePrice: dto.OptionalFloat{
			Set:   true,
			Value: &newSalePrice,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "sale_price must be lower than price", err.Error())
}

func TestPatchProductService_SalePriceGreaterThanZero_ShouldReturnBadRequest(t *testing.T) {
	id := "TESTPATCHPRODUCT"
	salePrice := float64(1278)

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)
			return &domain.Product{
				ID:        id,
				Name:      "Johnson's Baby",
				SalePrice: &salePrice,
				Price:     float64(1680),
			}, nil
		},
	}
	svc := NewProductService(repo)

	newSalePrice := float64(-1500)
	req := dto.RequestPatchProduct{
		SalePrice: dto.OptionalFloat{
			Set:   true,
			Value: &newSalePrice,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "sale_price must be greater than or equal 0", err.Error())
}

func TestPatchProductService_PriceGreaterThanZero_ShouldReturnBadRequest(t *testing.T) {
	id := "TESTPATCHPRODUCT"
	salePrice := float64(1278)

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)
			return &domain.Product{
				ID:        id,
				Name:      "Johnson's Baby",
				SalePrice: &salePrice,
				Price:     float64(1680),
			}, nil
		},
	}
	svc := NewProductService(repo)

	newPrice := float64(-1680)
	req := dto.RequestPatchProduct{
		Price: dto.OptionalFloat{
			Set:   true,
			Value: &newPrice,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "price must be greater than or equal 0", err.Error())
}

func TestPatchProductService_UpdateRepositoryError_ShouldReturnInternalServerError(t *testing.T) {
	id := "TESTPATCHPRODUCT"

	repo := &mockProductRepository{
		GetProductByIDFunc: func(ctx context.Context, productID string) (*domain.Product, error) {
			assert.Equal(t, id, productID)
			return &domain.Product{
				ID:    id,
				Name:  "Johnson's Baby",
				Price: float64(1680),
			}, nil
		},
		UpdateFunc: func(ctx context.Context, productID string, req *dto.RequestPatchProduct) error {
			return errors.New("server error")
		},
	}
	svc := NewProductService(repo)

	newName := "New Name"
	req := dto.RequestPatchProduct{
		Name: dto.OptionalString{
			Set:   true,
			Value: &newName,
		},
	}

	status, err := svc.PatchProductService(context.Background(), id, &req)
	require.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, "server error", err.Error())
}
