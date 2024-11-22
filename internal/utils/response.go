package utils

import (
	"ByteScience-WAM-Business/internal/model/dto"
	"github.com/gin-gonic/gin"
)

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) dto.Response {
	return dto.Response{
		Code:    Success,
		Message: ErrorMessages[Success],
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string) dto.Response {
	// 如果没有提供自定义消息，则使用错误码对应的消息
	if message == "" {
		message = ErrorMessages[code]
	}

	return dto.Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// SendResponse 发送响应
func SendResponse(ctx *gin.Context, statusCode int, response dto.Response) {
	ctx.JSON(statusCode, response)
	ctx.Abort() // 中断后续处理逻辑，确保响应后停止执行
}
