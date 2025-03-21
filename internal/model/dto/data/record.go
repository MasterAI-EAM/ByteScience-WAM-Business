package data

import "ByteScience-WAM-Business/internal/model/dto"

// RecordListRequest 分页查询配方列表的请求体结构
type RecordListRequest struct {
	// PaginationRequest 通用分页参数，包括页码和每页大小等信息
	dto.PaginationRequest

	// RecordName 配方名称（可选），用于按名称前缀模糊匹配查询
	// 约束：长度 2-255 个字符
	RecordName string `json:"recordName" validate:"omitempty,min=2,max=255" example:"配方名称"`

	// OpType 操作类型（可选），用于按操作类型筛选配方
	// 约束：长度 2-255 个字符
	OpType string `json:"opType" validate:"omitempty,min=2,max=255" example:"查询"`

	// UserId 操作用户的 ID，UUID 格式
	UserId string `json:"userId" example:"f8792732-9740-47fd-ae87-a4795e3d2045"`
}

// RecordListResponse 分页查询配方列表的响应体结构
type RecordListResponse struct {
	// Total 符合条件的配方总数
	Total int64 `json:"total" example:"100"`

	// List 当前页返回的配方列表
	List []RecordData `json:"list"`
}

// RecordData 配方记录详情
type RecordData struct {
	// RecordId 配方 ID，唯一标识，UUID 格式
	RecordId string `json:"recordId" example:"f8792732-9740-47fd-ae87-a4795e3d2045"`

	// OpTypeName 操作类型名称，例如："登录"
	OpTypeName string `json:"opTypeName" example:"登录"`

	// OpType 操作类型的数值表示，例如：1 代表登录操作
	OpType string `json:"opType" example:"login"`

	// UserName 执行操作的用户姓名
	UserName string `json:"userName" example:"小明"`

	// UserId 操作用户的 ID，UUID 格式
	UserId string `json:"userId" example:"f8792732-9740-47fd-ae87-a4795e3d2045"`

	// Desc 操作描述，例如："用户登录系统"
	Desc string `json:"desc" example:"用户登录系统"`

	// CreatedAt 记录创建时间，ISO 8601 格式
	CreatedAt string `json:"createdAt" example:"2006-01-02 15:04:05"`
}
