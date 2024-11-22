package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// bindAndValidate 处理请求的参数绑定和校验
func bindAndValidate[T any](ctx *gin.Context, method string) (*T, error) {
	var req T
	var err error

	err = ctx.ShouldBindJSON(&req)

	// 请求参数绑定失败
	if err != nil {
		return nil, fmt.Errorf("Invalid Request Parameters: %v", err)
	}

	// 参数校验
	validate := validator.New()

	// 执行校验
	err = validate.Struct(req)
	if err != nil {
		// 处理校验错误，返回最关键的错误信息
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			// 获取字段名和具体的校验错误信息
			errorMessage := fmt.Sprintf("Field '%s' %s", e.Field(), e.Tag())

			// 检查具体的值和客户端传入的值
			if e.Param() != "" {
				// 传入的值和期望的值
				errorMessage = fmt.Sprintf("%s (expected: %s)", errorMessage, e.Param())
			}

			validationErrors = append(validationErrors, errorMessage)
		}

		// 如果有错误，返回第一个错误信息
		if len(validationErrors) > 0 {
			return nil, fmt.Errorf(validationErrors[0])
		}
		return nil, fmt.Errorf("Invalid Request Parameters")
	}

	return &req, nil
}
