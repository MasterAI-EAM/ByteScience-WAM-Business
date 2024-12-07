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

// PaginationRequest 通用分页请求
type PaginationRequest struct {
	// Page 页码，选填，范围限制：[1,10000]
	// 用于分页查询管理员列表，最小值为1，最大值为10000
	Page int `json:"page" validate:"omitempty,gte=1,lte=10000" example:"1"`

	// PageSize 每页大小，选填，范围限制：[1,10000]
	// 用于限制每页返回的管理员数量，最小值为1，最大值为10000
	PageSize int `json:"pageSize" validate:"omitempty,gte=1,lte=10000" example:"10"`
}
