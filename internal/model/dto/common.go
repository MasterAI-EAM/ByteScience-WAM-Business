package dto

// Empty 空结构体
type Empty struct{}

// Response 成功响应格式
type Response struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 信息
	Data    interface{} `json:"data"`    // 响应数据
}

// ErrorResponse 错误响应格式
type ErrorResponse struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}
