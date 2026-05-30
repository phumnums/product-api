package service

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"main.go/internal/domain"
	"main.go/internal/dto"
	"main.go/internal/mapper"
	"main.go/internal/repository"
)

type ProductUsecase interface {
	CreateProductService(ctx context.Context, req dto.RequestCreateProduct) (*dto.ResponseCreateProduct, int, error)
}

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (u *ProductService) CreateProductService(ctx context.Context, req dto.RequestCreateProduct) (*dto.ResponseCreateProduct, int, error) {

	if req.SalePrice != nil && *req.SalePrice >= req.Price {
		return nil, http.StatusBadRequest, errors.New("sale_price must be lower than price")
	}

	productData := domain.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	product, err := u.productRepo.Create(ctx, productData)
	if err != nil {
		log.Printf("[CreateProductService] create product repo err -> %s", err)
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	return mapper.ProductResponseMapper(product), http.StatusOK, nil

}
