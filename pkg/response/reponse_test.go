package response

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCreateProductRequest struct {
	Name      string   `validate:"required"`
	Price     *float64 `validate:"required,gte=0"`
	SalePrice *float64 `validate:"omitempty,gte=0"`
}

func TestVAlidateionError_RequiredField_ShoudReturnMessages(t *testing.T) {
	validate := validator.New()
	req := testCreateProductRequest{}

	err := validate.Struct(req)
	require.Error(t, err)

	validationErrors, ok := err.(validator.ValidationErrors)
	require.True(t, ok)

	result := ValidationError(validationErrors)

	assert.Contains(t, result, "Name is required")
	assert.Contains(t, result, "Price is required")

}

func TestValidationError_GTE_ShoudReturnMessage(t *testing.T) {
	validate := validator.New()

	salePrice := float64(-10)
	price := float64(-20)

	req := testCreateProductRequest{
		Name:      "Johnson's Baby",
		SalePrice: &salePrice,
		Price:     &price,
	}

	err := validate.Struct(req)
	require.Error(t, err)

	validationErrors, ok := err.(validator.ValidationErrors)
	require.True(t, ok)

	result := ValidationError(validationErrors)

	assert.Contains(t, result, "SalePrice must be greater than or equal 0")
	assert.Contains(t, result, "Price must be greater than or equal 0")
}
