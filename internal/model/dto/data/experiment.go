package data

import "ByteScience-WAM-Business/internal/model/dto"

// ExperimentListRequest 用于分页查询实验列表的请求体结构
type ExperimentListRequest struct {
	// PaginationRequest dto.PaginationRequest 通用分页参数
	// 包含页码和每页大小等分页信息
	dto.PaginationRequest

	// ExperimentName string 实验名称，选填，长度限制：2-128字符
	// 用于按名称模糊查询实验记录
	ExperimentName string `json:"experimentName" validate:"omitempty,min=2,max=128" example:"实验名称"`

	// Experimenter string 实验者，选填，长度限制：2-128字符
	// 进行实验的人员名称
	Experimenter string `json:"experimenter" validate:"omitempty,min=1,max=128" example:"张三"`
}

// ExperimentListResponse 分页查询实验列表的响应体结构
type ExperimentListResponse struct {
	// Total int64 总条数
	// 返回符合条件的实验记录总数
	Total int64 `json:"total" example:"100"`

	// List []ExperimentData 数据
	// 分页返回的实验记录列表
	List []ExperimentData `json:"list"`
}

// ExperimentData 实验数据结构
type ExperimentData struct {
	// ExperimentID string 实验ID
	// 唯一标识实验的UUID
	ExperimentID string `json:"experimentId" example:"123e4567-e89b-12d3-a456-426614174000"`

	// ExperimentName string 实验名称
	// 实验的名称信息
	ExperimentName string `json:"experimentName" example:"实验名称"`

	// FileID string 文件ID
	// 关联的文件资源ID
	FileID string `json:"fileId" example:"123e4567-e89b-12d3-a456-426614174001"`

	// fileName string 文件名
	// 关联的文件资源ID
	FileName string `json:"fileName" example:"240628AI模型数据200组 含FRP性能-(对外）FD"`

	// EntryCategory int8 录入类别
	// 1 表示文件导入，2 表示页面输入
	EntryCategory int8 `json:"entryCategory" example:"1"`

	// Experimenter string 实验者
	// 进行实验的人员名称
	Experimenter string `json:"experimenter" example:"张三"`

	// UserID string 操作用户ID
	// 记录操作该实验的用户 ID
	UserID string `json:"userId" example:"987e6543-d21b-34c5-a654-123456789abc"`

	// StartTime 实验开始时间
	// 格式为时间戳，实验开始时间
	StartTime string `json:"startTime" example:"2024-11-18T10:00:00Z"`

	// EndTime 实验结束时间
	// 格式为时间戳，实验结束时间
	EndTime string `json:"endTime" example:"2024-11-18T10:00:00Z"`

	// CreatedAt 创建时间
	// 格式为时间戳，创建时间
	CreatedAt string `json:"createdAt" example:"2024-11-18T10:00:00Z"`

	// Steps []ExperimentStepData 实验步骤
	// 包含该实验的步骤信息
	Steps []ExperimentStepData `json:"steps"`
}

// ExperimentStepData 实验步骤数据结构
type ExperimentStepData struct {
	// StepID string 步骤ID
	// 唯一标识实验步骤的UUID
	StepID string `json:"stepId" example:"123e4567-e89b-12d3-a456-426614174002"`

	// RecipeID string 配方ID
	// 配方的UUID
	RecipeID string `json:"recipeId" example:"123e4567-e89b-12d3-a456-426614174002"`

	// StepName string 步骤名称
	// 描述实验步骤的名称
	StepName string `json:"stepName" example:"步骤名称"`

	// StepNameDescription string 实验步骤描述
	// 实验步骤描述
	StepNameDescription string `json:"stepNameDescription" example:"实验步骤描述"`

	// ExperimentCondition string 实验条件
	// 步骤对应的实验条件描述
	ExperimentCondition string `json:"experimentCondition" example:"实验条件"`

	// StepOrder int 排序(从大到小)
	// 实验步骤的执行排序(从大到小)
	// 示例值: 1
	StepOrder int `json:"stepOrder" validate:"required,gte=1" example:"1"`

	// ResultValue string 实验条件
	// 步骤对应的实验结果
	ResultValue string `json:"resultValue" example:"步骤结果值"`
}

// ExperimentDeleteRequest 删除实验
type ExperimentDeleteRequest struct {
	// experimentId 实验编号，必填，UUID格式
	// 唯一标识要删除的实验，格式必须为UUID4
	ExperimentID string `json:"experimentId" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// ExperimentUpdateRequest 修改实验数据
type ExperimentUpdateRequest struct {
	// ExperimentID string 实验ID
	// 唯一标识实验的 UUID，必填，用于确定要更新的实验记录
	// 示例值: "123e4567-e89b-12d3-a456-426614174000"
	ExperimentID string `json:"experimentId" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`

	// ExperimentName string 实验名称
	// 实验的名称，选填；如果填写，名称长度限制为 2-128 字符
	// 支持中文、英文、数字及特殊字符
	// 示例值: "实验名称"
	ExperimentName string `json:"experimentName" validate:"omitempty,min=2,max=128" example:"实验名称"`

	// Experimenter string 实验者
	// 实验的负责人，选填；如果填写，名称长度限制为 1-128 字符
	// 示例值: "张三"
	Experimenter string `json:"experimenter" validate:"omitempty,min=1,max=128" example:"张三"`

	// StartTime string 实验开始时间
	// 选填，格式为 "2006-01-02T15:04:05Z"（RFC3339 格式）
	// 示例值: "2024-02-05T08:30:00Z"
	StartTime string `json:"startTime" validate:"omitempty,datetime=2006-01-02T15:04:05Z" example:"2024-02-05T08:30:00Z"`

	// EndTime string 实验结束时间
	// 选填，格式为 "2006-01-02T15:04:05Z"（RFC3339 格式）
	// 示例值: "2024-02-05T18:00:00Z"
	EndTime string `json:"endTime" validate:"omitempty,datetime=2006-01-02T15:04:05Z" example:"2024-02-05T18:00:00Z"`

	// Steps []ExperimentStepUpdate 实验步骤列表
	// 选填，包含该实验的所有步骤信息，每个步骤包含名称、描述、实验条件、结果值及材料组等
	// 示例值: [{"stepName": "步骤名称", "stepNameDescription": "实验步骤描述", "experimentCondition": "实验条件", "resultValue": "步骤结果值", "materialGroups": []}]
	Steps []ExperimentStepUpdate `json:"steps" validate:"omitempty,dive"`
}

// ExperimentStepUpdate 修改实验步骤请求结构
type ExperimentStepUpdate struct {
	// StepName string 步骤名称
	// 实验步骤的名称，必填，限制长度为 1-255 字符
	// 示例值: "步骤名称"
	StepName string `json:"stepName" validate:"required,min=1,max=255" example:"步骤名称"`

	// ExperimentCondition string 实验条件
	// 选填，实验步骤的实验条件描述，最长 255 字符
	// 示例值: "实验条件"
	ExperimentCondition string `json:"experimentCondition" validate:"omitempty,max=255" example:"实验条件"`

	// ResultValue string 步骤结果值
	// 选填，实验步骤的结果值，最长 256 字符
	// 示例值: "步骤结果值"
	ResultValue string `json:"resultValue" validate:"required,max=256" example:"步骤结果值"`

	// StepOrder int 排序(从大到小)
	// 必填，实验步骤的执行排序(从大到小)，必须为正整数
	// 示例值: 1
	StepOrder int `json:"stepOrder" validate:"required,gte=1" example:"1"`

	// RecipeID string 配方ID
	// 配方的UUID
	RecipeID string `json:"recipeId" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174002"`
}

// ExperimentAddRequest 添加实验数据
type ExperimentAddRequest struct {
	// ExperimentName string 实验名称
	// 实验的名称，选填；如果填写，名称长度限制为 2-128 字符
	// 支持中文、英文、数字及特殊字符
	// 示例值: "实验名称"
	ExperimentName string `json:"experimentName" validate:"omitempty,min=2,max=128" example:"实验名称"`

	// Experimenter string 实验者
	// 实验的负责人，选填；如果填写，名称长度限制为 1-128 字符
	// 示例值: "张三"
	Experimenter string `json:"experimenter" validate:"omitempty,min=1,max=128" example:"张三"`

	// Sort int 排序 优先级从大到小
	// 选填，必须是大于等于 0 的整数
	// 示例值: 1
	Sort int `json:"sort" validate:"omitempty,min=0" example:"1"`

	// StartTime string 实验开始时间
	// 选填，格式为 "2006-01-02T15:04:05Z"（RFC3339 格式）
	// 示例值: "2024-02-05T08:30:00Z"
	StartTime string `json:"startTime" validate:"omitempty,datetime=2006-01-02T15:04:05Z" example:"2024-02-05T08:30:00Z"`

	// EndTime string 实验结束时间
	// 选填，格式为 "2006-01-02T15:04:05Z"（RFC3339 格式）
	// 示例值: "2024-02-05T18:00:00Z"
	EndTime string `json:"endTime" validate:"omitempty,datetime=2006-01-02T15:04:05Z" example:"2024-02-05T18:00:00Z"`

	// Steps []ExperimentStepAdd 实验步骤列表
	// 选填，包含该实验的所有步骤信息，每个步骤包含名称、描述、实验条件、结果值及材料组等
	// 示例值: [{"stepName": "步骤名称", "stepNameDescription": "实验步骤描述", "experimentCondition": "实验条件", "resultValue": "步骤结果值", "materialGroups": []}]
	Steps []ExperimentStepAdd `json:"steps" validate:"omitempty,dive"`
}

// ExperimentStepAdd 添加实验步骤请求结构
type ExperimentStepAdd struct {
	// StepName string 步骤名称
	// 实验步骤的名称，必填，限制长度为 1-255 字符
	// 示例值: "步骤名称"
	StepName string `json:"stepName" validate:"required,min=1,max=255" example:"步骤名称"`

	// ExperimentCondition string 实验条件
	// 选填，实验步骤的实验条件描述，最长 255 字符
	// 示例值: "实验条件"
	ExperimentCondition string `json:"experimentCondition" validate:"omitempty,max=255" example:"实验条件"`

	// ResultValue string 步骤结果值
	// 选填，实验步骤的结果值，最长 256 字符
	// 示例值: "步骤结果值"
	ResultValue string `json:"resultValue" validate:"required,max=256" example:"步骤结果值"`

	// StepOrder int 排序(从大到小)
	// 必填，实验步骤的执行排序(从大到小)，必须为正整数
	// 示例值: 1
	StepOrder int `json:"stepOrder" validate:"required,gte=1" example:"1"`

	// RecipeID string 配方ID
	// 配方的UUID
	RecipeID string `json:"recipeId" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174002"`
}

// ExperimentStepInfo 实验步骤数据结构
type ExperimentStepInfo struct {
	// StepID string 步骤ID
	// 唯一标识实验步骤的 UUID，必填；用于确定需要更新的实验步骤记录
	// 示例值: "123e4567-e89b-12d3-a456-426614174002"
	StepID string `json:"stepId" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174002"`

	// RecipeID string 配方ID
	// 配方的UUID
	RecipeID string `json:"recipeId" example:"123e4567-e89b-12d3-a456-426614174002"`

	// StepName string 步骤名称
	// 描述实验步骤的名称，必填；长度限制为 1-255 字符
	// 示例值: "步骤名称"
	StepName string `json:"stepName" validate:"required,min=1,max=255" example:"步骤名称"`

	// StepNameDescription string 实验步骤描述
	// 对实验步骤的详细说明，选填；长度限制为 0-512 字符
	// 示例值: "该步骤描述了实验具体操作流程"
	StepNameDescription string `json:"stepNameDescription" validate:"omitempty,max=512" example:"实验步骤描述"`

	// ExperimentCondition string 实验条件
	// 用于描述实验步骤的前置条件或约束条件，选填；长度限制为 0-255 字符
	// 示例值: "需要在 25 摄氏度环境中进行"
	ExperimentCondition string `json:"experimentCondition" validate:"omitempty,max=255" example:"实验条件"`

	// ResultValue string 实验结果值
	// 记录实验步骤产生的结果数据，选填；长度限制为 0-256 字符
	// 示例值: "实验结果为成功"
	ResultValue string `json:"resultValue" validate:"omitempty,max=256" example:"步骤结果值"`
}
