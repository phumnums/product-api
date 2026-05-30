package mapper

import (
	"main.go/internal/domain"
	"main.go/internal/dto"
)

func ProductResponseMapper(product *domain.Product) *dto.ResponseCreateProduct {

	result := dto.ResponseCreateProduct{
		Name:        product.Name,
		Description: product.Description,
		SalePrice:   product.SalePrice,
		Price:       product.Price,
	}

	return &result
}
