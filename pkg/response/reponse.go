package response

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Successful bool        `json:"successful"`
	ErrorCode  string      `json:"error_code"`
	Data       interface{} `json:"data"`
}

type ResponseNoData struct {
	Successful bool   `json:"successful"`
	ErrorCode  string `json:"error_code"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Successful: true,
		ErrorCode:  "",
		Data:       data,
	})
}

func Error(c *gin.Context, statusCode int, errorCode string, data interface{}) {
	c.JSON(statusCode, Response{
		Successful: false,
		ErrorCode:  errorCode,
		Data:       data,
	})
}

func SuccessNoData(c *gin.Context, statusCode int) {
	c.JSON(statusCode, ResponseNoData{
		Successful: true,
		ErrorCode:  "",
	})
}

func ErrorNoData(c *gin.Context, statusCode int, errorCode string) {
	c.JSON(statusCode, ResponseNoData{
		Successful: false,
		ErrorCode:  errorCode,
	})
}

func ValidationError(errs validator.ValidationErrors) []string {
	var result []string

	for _, errField := range errs {
		field := errField.Field()

		switch errField.Tag() {
		case "required":
			result = append(result, fmt.Sprintf("%s is required", field))
		case "gte":
			result = append(result, fmt.Sprintf("%s must be greater than or equal %s", field, errField.Param()))
		default:
			result = append(result, fmt.Sprintf("%s is invalid", field))
		}
	}
	return result
}
