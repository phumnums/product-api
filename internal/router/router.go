package router

import (
	"github.com/gin-gonic/gin"
	"main.go/internal/handler"
)

func NewRouter(productHandler *handler.ProductHandler) *gin.Engine {
	r := gin.Default()

	registerProductRoutes(r, productHandler)

	return r
}

func registerProductRoutes(r *gin.Engine, productHandler *handler.ProductHandler) {

	product := r.Group("/product")
	{
		product.POST("", productHandler.CreateProductHandler)
	}
}
