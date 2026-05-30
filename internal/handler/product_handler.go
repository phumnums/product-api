package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/dto"
	service "main.go/internal/service"
	"main.go/pkg/response"
)

type ProductHandler struct {
	productUsecase service.ProductUsecase
}

func NewProductHandler(productUsecase service.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

func (h *ProductHandler) CreateProductHandler(c *gin.Context) {
	var req dto.RequestCreateProduct

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("err: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid request", response.ValidationError(err))
		return
	}

	product, status, err := h.productUsecase.CreateProductService(c, req)
	if err != nil {
		response.Error(c, status, err.Error(), nil)
		return
	}

	response.Success(c, status, product)
}
