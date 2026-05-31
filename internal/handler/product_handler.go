package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		if validationError, ok := err.(validator.ValidationErrors); ok {
			response.Error(c, http.StatusBadRequest, "invalid request", response.ValidationError(validationError))
			return
		}

		response.Error(c, http.StatusBadRequest, "invalid request", nil)
		return
	}

	product, status, err := h.productUsecase.CreateProductService(c, &req)
	if err != nil {
		response.Error(c, status, err.Error(), nil)
		return
	}

	response.Success(c, status, product)
}

func (h *ProductHandler) PatchProductHandler(c *gin.Context) {
	productID := c.Param("id")
	var req dto.RequestPatchProduct

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"successful": false,
			"error_code": err.Error(),
		})
		return
	}

	status, err := h.productUsecase.PatchProductService(c, productID, &req)
	if err != nil {
		c.JSON(status, gin.H{
			"successful": false,
			"error_code": err.Error(),
		})
		return
	}

	c.JSON(status, gin.H{
		"successful": true,
		"error_code": nil,
	})

}
