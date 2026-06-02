package component

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"main.go/internal/domain"
	"main.go/internal/handler"
	"main.go/internal/repository"
	"main.go/internal/router"
	"main.go/internal/service"
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

func setupTestRouter(t *testing.T) http.Handler {
	db := setupTestDB(t)

	productRepo := repository.NewProductPostgresRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	return router.NewRouter(productHandler)
}

func TestCreateProductComponent_Success(t *testing.T) {
	r := setupTestRouter(t)

	body := `{
    	"name": "Johnson's Baby",
		"description": "Top-To-Toe Hair&Body Bath 100 ml. x48",
		"sale_price": 1278,
		"price": 1680
	}`

	req := httptest.NewRequest(
		http.MethodPost, "/product",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), `"successful":true`)
	assert.Contains(t, res.Body.String(), `"error_code":""`)
	assert.Contains(t, res.Body.String(), `"name":"Johnson's Baby"`)
}

func TestCreateProductComponent_InvalidRequest(t *testing.T) {
	r := setupTestRouter(t)

	body := `{}`

	req := httptest.NewRequest(
		http.MethodPost, "/product",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Contains(t, res.Body.String(), `"successful":false`)
	assert.Contains(t, res.Body.String(), `"error_code":"invalid request"`)
}

func TestPatchProductComponent_Success(t *testing.T) {
	r := setupTestRouter(t)
	db := setupTestDB(t)
	repo := repository.NewProductPostgresRepository(db)

	product := &domain.Product{
		ID:    ulid.Make().String(),
		Name:  "Johnson's Baby",
		Price: float64(1680),
	}

	created, err := repo.Create(context.Background(), product)
	require.NoError(t, err)

	body := `{
		"name": "Johnson's Baby",
		"price": 1700
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/product/"+created.ID,
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), `"successful":true`)
	assert.Contains(t, res.Body.String(), `"error_code":""`)

	updated, err := repo.GetProductByID(context.Background(), created.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Johnson's Baby", updated.Name)
	assert.Equal(t, float64(1700), updated.Price)
}

func TestPatchProductComponent_NotFound(t *testing.T) {
	r := setupTestRouter(t)

	body := `{
		"name": "Johnson's Baby",
		"price": 1700
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/product/TESTPRODUCTNOTFOUND",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Contains(t, res.Body.String(), `"successful":false`)
	assert.Contains(t, res.Body.String(), `"error_code":"product not found"`)
}

func TestPatchProductComponent_EmptyBody(t *testing.T) {
	r := setupTestRouter(t)

	db := setupTestDB(t)
	repo := repository.NewProductPostgresRepository(db)

	product := &domain.Product{
		ID:    ulid.Make().String(),
		Name:  "Johnson's Baby",
		Price: float64(1680),
	}

	created, err := repo.Create(context.Background(), product)
	require.NoError(t, err)

	body := `{}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/product/"+created.ID,
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Contains(t, res.Body.String(), `"successful":false`)
	assert.Contains(t, res.Body.String(), `"error_code":"invalid request"`)
}
