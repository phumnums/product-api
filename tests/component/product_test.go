package component

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
