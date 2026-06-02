package dto

type RequestCreateProduct struct {
	Name        string   `json:"name" binding:"required" example:"Johnson's Baby"`
	Description *string  `json:"description" binding:"omitempty" example:"Top-To-Toe Hair&Body Bath"`
	SalePrice   *float64 `json:"sale_price" binding:"omitempty,gte=0" example:"1278"`
	Price       float64  `json:"price" binding:"required,gte=0" example:"1680"`
}

type ResponseCreateProduct struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}
type RequestPatchProduct struct {
	Name        OptionalString `json:"name"`
	Description OptionalString `json:"description"`
	SalePrice   OptionalFloat  `json:"sale_price"`
	Price       OptionalFloat  `json:"price"`
}
