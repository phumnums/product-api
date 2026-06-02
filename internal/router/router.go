package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "main.go/docs"
	"main.go/internal/handler"
)

func NewRouter(productHandler *handler.ProductHandler) *gin.Engine {
	r := gin.Default()

	registerProductRoutes(r, productHandler)
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func registerProductRoutes(r *gin.Engine, productHandler *handler.ProductHandler) {

	product := r.Group("/product")
	{
		product.POST("", productHandler.CreateProductHandler)
		product.PATCH("/:id", productHandler.PatchProductHandler)
	}
}
