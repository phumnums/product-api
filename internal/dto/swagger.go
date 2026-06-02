package dto

type RequestPatchProductSwagger struct {
	Name        *string  `json:"name" example:"New Product Name"`
	Description *string  `json:"description" example:"New Description"`
	SalePrice   *float64 `json:"sale_price" example:"1200"`
	Price       *float64 `json:"price" example:"1700"`
}

type PatchProductSuccessResponseSwagger struct {
	Successful bool   `json:"successful" example:"true"`
	ErrorCode  string `json:"error_code" example:""`
}

type PatchProductBadRequestResponseSwagger struct {
	Successful bool   `json:"successful" example:"false"`
	ErrorCode  string `json:"error_code" example:"invalid request"`
}

type PatchProductNotFoundResponseSwagger struct {
	Successful bool   `json:"successful" example:"false"`
	ErrorCode  string `json:"error_code" example:"product not found"`
}

type PatchProductServerErrorResponseSwagger struct {
	Successful bool   `json:"successful" example:"false"`
	ErrorCode  string `json:"error_code" example:"server error"`
}

type CreateProductSuccessResponseSwagger struct {
	Successful bool   `json:"successful" example:"true"`
	ErrorCode  string `json:"error_code" example:""`
	Data       struct {
		Name        string   `json:"name" example:"Johnson's Baby"`
		Description *string  `json:"description" example:"Top-To-Toe Hair&Body Bath"`
		SalePrice   *float64 `json:"sale_price" example:"1278"`
		Price       float64  `json:"price" example:"1680"`
	} `json:"data"`
}

type CreateProductValodationErrorSwagger struct {
	Successful bool     `json:"successful" example:"false"`
	ErrorCode  string   `json:"error_code" example:"invalid request"`
	Data       []string `json:"data" example:"Name is required"`
}

type CreateProductBadRequestSwagger struct {
	Successful bool    `json:"successful" example:"false"`
	ErrorCode  string  `json:"error_code" example:"invalid request"`
	Data       *string `json:"data" example:"null"`
}

type CreateProductServerErrorSwagger struct {
	Successful bool    `json:"successful" example:"false"`
	ErrorCode  string  `json:"error_code" example:"server error"`
	Data       *string `json:"data" example:"null"`
}
