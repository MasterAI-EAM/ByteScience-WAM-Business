// utils/business_error.go

package utils

import "fmt"

// BusinessError 业务错误类型
type BusinessError struct {
	Code    int
	Message string
}

// 实现 error 接口
func (e *BusinessError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewBusinessError 创建一个新的业务错误
func NewBusinessError(code int) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: ErrorMessages[code],
	}
}
