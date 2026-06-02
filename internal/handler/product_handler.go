package handler

import (
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

// CreateProductHandler godoc
//
// @Summary Create Product
// @Description Create new product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.RequestCreateProduct true "Create Product"
// @Success 200 {object} dto.CreateProductSuccessResponseSwagger
// @Failure 400 {object} dto.CreateProductValodationErrorSwagger
// @Failure 500 {object} dto.CreateProductServerErrorSwagger
// @Router /product [post]
func (h *ProductHandler) CreateProductHandler(c *gin.Context) {
	var req dto.RequestCreateProduct

	if err := c.ShouldBindJSON(&req); err != nil {
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

// PatchProductHandler godoc
//
// @Summary Patch Product
// @Description Update product partially
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.RequestPatchProductSwagger true "Patch Product"
// @Success 200 {object} dto.PatchProductSuccessResponseSwagger
// @Failure 400 {object} dto.PatchProductBadRequestResponseSwagger
// @Failure 404 {object} dto.PatchProductNotFoundResponseSwagger
// @Failure 500 {object} dto.PatchProductServerErrorResponseSwagger
// @Router /product/{id} [patch]
func (h *ProductHandler) PatchProductHandler(c *gin.Context) {
	productID := c.Param("id")
	var req dto.RequestPatchProduct

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorNoData(c, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.productUsecase.PatchProductService(c, productID, &req)
	if err != nil {
		response.ErrorNoData(c, status, err.Error())
		return
	}
	response.SuccessNoData(c, status)
}
