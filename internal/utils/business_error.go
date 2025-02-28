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

// NewBusinessError 创建一个新的业务错误实例
func NewBusinessError(code int, message string, args ...interface{}) *BusinessError {
	if message == "" {
		// 如果没有传入 message，使用默认的 error template
		messageTemplate, exists := ErrorMessages[code]
		if !exists {
			// 如果找不到对应的错误信息模板，使用默认模板
			messageTemplate = "Unknown error code: %d"
			args = []interface{}{code}
		}
		// 使用 fmt.Sprintf 对错误信息模板进行格式化
		message = fmt.Sprintf(messageTemplate, args...)
	}

	return &BusinessError{
		Code:    code,
		Message: message,
	}
}
