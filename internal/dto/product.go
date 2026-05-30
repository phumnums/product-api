package dto

type RequestCreateProduct struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description" binding:"omitempty"`
	SalePrice   *float64 `json:"sale_price" binding:"omitempty,gte=0"`
	Price       float64  `json:"price" binding:"required,gte=0"`
}

type ResponseCreateProduct struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}
