package data

import "ByteScience-WAM-Business/internal/model/dto"

// TaskListRequest 用于分页查询任务列表的请求体结构
type TaskListRequest struct {
	// PaginationRequest dto.PaginationRequest 通用分页参数
	// 包含页码和每页大小等分页信息
	dto.PaginationRequest

	// FileName string 文件名称，选填，长度限制：2-128字符
	// 用于按文件名称模糊查询任务记录
	FileName string `json:"fileName" validate:"omitempty,min=2,max=128" example:"文件名称"`

	// Status string 任务状态
	// 当前任务的处理状态，pending=待处理, processing=处理中, success=成功, failure=失败
	Status string `json:"status" validate:"omitempty,oneof=pending processing success failure" example:"pending"`

	// StartTime string 开始时间
	// 任务的创建时间，格式 2006-01-02 15:04:05
	StartTime string `json:"startTime" validate:"omitempty,datetime=2006-01-02 15:04:05" example:"2006-01-02 15:04:05"`

	// EndTime string 结束时间
	// 任务的创建时间，格式 2006-01-02 15:04:05
	EndTime string `json:"endTime" validate:"omitempty,datetime=2006-01-02 15:04:05" example:"2006-01-02 15:04:05"`
}

// TaskListResponse 分页查询任务列表的响应体结构
type TaskListResponse struct {
	// Total int64 总条数
	// 返回符合条件的任务记录总数
	Total int64 `json:"total" example:"100"`

	// List []TaskData 数据
	// 分页返回的任务记录列表
	List []TaskData `json:"list"`
}

// TaskData 任务数据结构
type TaskData struct {
	// TaskID string 任务ID
	// 唯一标识任务的UUID
	TaskID string `json:"taskId" example:"123e4567-e89b-12d3-a456-426614174000"`

	// BatchID string 批次号
	// 用于标识一组相关任务的批次UUID
	BatchID string `json:"batchId" example:"987e6543-e89b-12d3-a456-426614174001"`

	// FileName string 文件名称
	// 任务关联的文件名称
	FileName string `json:"fileName" example:"240628AI模型数据200组 含FRP性能-(对外）FD"`

	// FilePath string 文件路径
	// 存储文件的路径，指向文件在服务器上的位置
	FilePath string `json:"filePath" example:"/uploads/2024/11/240628AI模型数据200组 含FRP性能-(对外）FD.json"`

	// JSONFilePath string 硬代码json文件路径
	// JSON 文件存储的路径
	JSONFilePath string `json:"jsonFilePath" example:"/uploads/2024/11/240628AI模型数据200组 含FRP性能-(对外）FD.json"`

	// AiFilePath string ai处理后json文件路径
	// AI 处理后的 JSON 文件存储路径
	AiFilePath string `json:"aiFilePath" example:"/uploads/2024/11/240628AI模型数据200组 含FRP性能-(对外）FD-ai.json"`

	// Status string 任务状态
	// 当前任务的处理状态，pending=待处理, processing=处理中, success=成功, failure=失败
	Status string `json:"status" example:"pending"`

	// Remark string 任务状态描述
	// 对任务状态的额外描述信息，如错误信息等
	Remark string `json:"remark" example:"任务正在处理中"`

	// CreatedAt string 创建时间
	// 任务的创建时间，格式 2006-01-02 15:04:05
	CreatedAt string `json:"createdAt" example:"2006-01-02 15:04:05"`

	// UpdatedAt string 修改时间
	// 任务的最后更新时间，采用 ISO 8601 格式（例如：2024-11-18T12:00:00Z）
	UpdatedAt string `json:"updatedAt" example:"2024-11-18T12:00:00Z"`
}
