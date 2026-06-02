// @title Product API
// @version 1.0
// @description Product Management API
// @host localhost:3000
// @BasePath /
package main

import (
	"context"
	"log"

	"main.go/config"
	"main.go/internal/handler"
	"main.go/internal/repository"
	"main.go/internal/router"
	service "main.go/internal/service"
	"main.go/pkg/database"
)

func main() {

	ctx := context.Background()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	productRepo := repository.NewProductPostgresRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r := router.NewRouter(productHandler)

	if err := r.Run(":3000"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
