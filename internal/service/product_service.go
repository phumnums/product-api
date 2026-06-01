package service

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"main.go/internal/domain"
	"main.go/internal/dto"
	"main.go/internal/mapper"
	"main.go/internal/repository"
)

type ProductUsecase interface {
	CreateProductService(ctx context.Context, req *dto.RequestCreateProduct) (*dto.ResponseCreateProduct, int, error)
	PatchProductService(ctx context.Context, id string, req *dto.RequestPatchProduct) (int, error)
}

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProductService(ctx context.Context, req *dto.RequestCreateProduct) (*dto.ResponseCreateProduct, int, error) {

	if req.SalePrice != nil && *req.SalePrice >= req.Price {
		return nil, http.StatusBadRequest, errors.New("sale_price must be lower than price")
	}

	// check blank and empty string
	if req.Description != nil {
		description := strings.TrimSpace(*req.Description)
		if description == "" {
			req.Description = nil
		}
	}

	productData := domain.Product{
		ID:          string(ulid.Make().String()),
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	product, err := s.productRepo.Create(ctx, &productData)
	if err != nil {
		log.Printf("[CreateProductService] create product repo err -> %s", err)
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	return mapper.ProductResponseMapper(product), http.StatusOK, nil

}

func (s *ProductService) PatchProductService(ctx context.Context, productID string, req *dto.RequestPatchProduct) (int, error) {

	// check product id
	if productID == "" {
		return http.StatusBadRequest, errors.New("product_id is required")
	}

	// check name for update
	if req.Name.Set {
		if req.Name.Value != nil {
			if strings.TrimSpace(*req.Name.Value) == "" {
				return http.StatusBadRequest, errors.New("name is required")
			}
		} else {
			return http.StatusBadRequest, errors.New("name is required")
		}
	}

	// check and get product by id
	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return http.StatusNotFound, errors.New("product not found")
		}
		log.Printf("[PatchProductService] get product by id err -> %s", err)
		return http.StatusInternalServerError, errors.New("server error")
	}

	// check data update
	if !req.Name.Set && !req.Description.Set && !req.SalePrice.Set && !req.Price.Set {
		return http.StatusBadRequest, errors.New("invalid request")
	}

	// check description for update
	if req.Description.Set {
		if req.Description.Value != nil {
			if strings.TrimSpace(*req.Description.Value) == "" {
				req.Description.Value = nil
			}
		} else {
			req.Description.Value = nil
		}
	}

	// check price for update
	oldPrice := product.Price
	if req.Price.Set {
		if req.Price.Value == nil {
			return http.StatusBadRequest, errors.New("price is required")
		}
		if *req.Price.Value < 0 {
			return http.StatusBadRequest, errors.New("price must be greater than or equal 0")
		}
		oldPrice = *req.Price.Value
	}

	// check sale_price for update
	if req.SalePrice.Set {
		if req.SalePrice.Value != nil {
			if *req.SalePrice.Value >= oldPrice {
				return http.StatusBadRequest, errors.New("sale_price must be lower than price")
			}
			if *req.SalePrice.Value < 0 {
				return http.StatusBadRequest, errors.New("sale_price must be greater than or equal 0")
			}
		}
	}

	// update data
	if err := s.productRepo.Update(ctx, productID, req); err != nil {
		log.Printf("[PatchProductService] update product repo err -> %s", err)
		return http.StatusInternalServerError, errors.New("server error")
	}

	return http.StatusOK, nil
}
