package data

import "ByteScience-WAM-Business/internal/model/dto"

// ExperimentListRequest 用于分页查询实验列表的请求体结构
type ExperimentListRequest struct {
	// PaginationRequest dto.PaginationRequest 通用分页参数
	// 包含页码和每页大小等分页信息
	dto.PaginationRequest

	// TaskId string 任务id
	// 唯一标识实验的UUID
	TaskId string `json:"taskId" validate:"omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`

	// ExperimentName string 实验名称，选填，长度限制：2-128字符
	// 用于按名称模糊查询实验记录
	ExperimentName string `json:"experimentName" validate:"omitempty,min=2,max=128" example:"实验名称"`

	// Experimenter string 实验者，选填，长度限制：2-128字符
	// 进行实验的人员名称
	Experimenter string `json:"experimenter" validate:"omitempty,min=1,max=128" example:"张三"`

	// Status string 审核状态
	// 当前任务的处理状态，pending_review=待审核, approved=审核通过, rejected=审核不通过
	Status string `json:"status" validate:"omitempty,oneof=pending_review approved rejected" example:"pending_review"`

	// StartTime string 开始时间
	// 创建时间，格式 2006-01-02 15:04:05
	StartTime string `json:"startTime" validate:"omitempty,datetime=2006-01-02 15:04:05" example:"2006-01-02 15:04:05"`

	// EndTime string 结束时间
	// 创建时间，格式 2006-01-02 15:04:05
	EndTime string `json:"endTime" validate:"omitempty,datetime=2006-01-02 15:04:05" example:"2006-01-02 15:04:05"`
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

	// EntryCategory string 录入类别
	// file_import=文件导入, manual_entry=页面录入
	EntryCategory string `json:"entryCategory" example:"1"`

	// Experimenter string 实验者
	// 进行实验的人员名称
	Experimenter string `json:"experimenter" example:"张三"`

	// UserID string 操作用户ID
	// 记录操作该实验的用户 ID
	UserID string `json:"userId" example:"987e6543-d21b-34c5-a654-123456789abc"`

	// Username 用户名
	// 表示实验相关的用户名
	Username string `json:"username" example:"testuser"`

	// Status 实验状态
	// 表示实验的审核状态，包含 'pending_review'、'approved'、'rejected'
	Status string `json:"status" example:"pending_review" enum:"pending_review,approved,rejected"`

	// StartTime 实验开始时间
	// 格式为时间戳，实验开始时间
	StartTime string `json:"startTime" example:"2006-01-02 15:04:05"`

	// EndTime 实验结束时间
	// 格式为时间戳，实验结束时间
	EndTime string `json:"endTime" example:"2006-01-02 15:04:05"`

	// CreatedAt 创建时间
	// 格式为时间戳，创建时间
	CreatedAt string `json:"createdAt" example:"2006-01-02 15:04:05"`
}

// ExperimentDeleteRequest 删除实验
type ExperimentDeleteRequest struct {
	// experimentId 实验编号，必填，UUID格式
	// 唯一标识要删除的实验，格式必须为UUID4
	ExperimentID string `json:"experimentId" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// ExperimentInfoRequest 实验详情请求结构
type ExperimentInfoRequest struct {
	// ExperimentID string 实验ID
	// 唯一标识实验的UUID
	ExperimentID string `json:"experimentId" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// ExperimentInfoResponse 实验详情
type ExperimentInfoResponse struct {
	// experimentId 实验编号，必填，UUID格式
	// 唯一标识要删除的实验，格式必须为UUID4
	ExperimentID string `json:"experimentId" example:"123e4567-e89b-12d3-a456-426614174000"`

	// ExperimentName string 实验名称
	// 实验的名称信息
	ExperimentName string `json:"experimentName" example:"实验名称"`

	// EntryCategory string 录入类别
	// file_import=文件导入, manual_entry=页面录入
	EntryCategory string `json:"entryCategory" example:"1"`

	// Experimenter string 实验者
	// 进行实验的人员名称
	Experimenter string `json:"experimenter" example:"张三"`

	// UserID string 操作用户ID
	// 记录操作该实验的用户 ID
	UserID string `json:"userId" example:"987e6543-d21b-34c5-a654-123456789abc"`

	// Username 用户名
	// 表示实验相关的用户名
	Username string `json:"username" example:"testuser"`

	// Status 实验状态
	// 表示实验的审核状态，包含 'pending_review'、'approved'、'rejected'
	Status string `json:"status" example:"pending_review" enum:"pending_review,approved,rejected"`

	// StartTime 实验开始时间
	// 格式为时间戳，实验开始时间
	StartTime string `json:"startTime" example:"2006-01-02 15:04:05"`

	// EndTime 实验结束时间
	// 格式为时间戳，实验结束时间
	EndTime string `json:"endTime" example:"2006-01-02 15:04:05"`

	// CreatedAt 创建时间
	// 格式为时间戳，创建时间
	CreatedAt string `json:"createdAt" example:"2006-01-02 15:04:05"`

	// StepInfo []ExperimentStepInfo 实验步骤数据
	// 该步骤中涉及的材料组信息
	StepInfo []ExperimentStepInfo `json:"stepInfo"`

	// MaterialGroups []MaterialGroupInfo 材料组
	// 该步骤中涉及的材料组信息
	MaterialGroups []MaterialGroupInfo `json:"materialGroups"`
}

// ExperimentStepInfo 实验步骤数据结构
type ExperimentStepInfo struct {
	// StepID string 步骤ID
	// 唯一标识实验步骤的UUID
	StepID string `json:"stepId" example:"123e4567-e89b-12d3-a456-426614174002"`

	// StepName string 步骤名称
	// 描述实验步骤的名称
	StepName string `json:"stepName" example:"步骤名称"`

	// StepCondition string 实验条件
	// 步骤对应的实验条件描述
	StepCondition string `json:"stepCondition" example:"步骤实验条件"`

	// StepCategory 实验步骤类别
	// 实验步骤的类型，包含 'resin_mixing'、'hardener_mixing'、'resin_hardener_mixing'、'mechanical_performance'
	StepCategory string `json:"stepCategory" example:"resin_mixing" enum:"resin_mixing,hardener_mixing,resin_hardener_mixing,mechanical_performance"`

	// ResultValue string 实验条件
	// 步骤对应的实验结果
	// resin_mixing=树脂混合 {"树脂粘度":{"温度":27,"粘度":1350},"环氧当量":26}
	// hardener_mixing=固化剂混合 {"胺值":9.5,"固化剂粘度":{"温度":27,"粘度":null}}
	// resin_hardener_mixing=树脂/固化剂混合 {"温度":"27℃","可用时间":140,"混合粘度":276}
	// mechanical_performance=力学性能 {"value": 79}
	ResultValue string `json:"resultValue" example:"步骤结果值"`
}

// MaterialGroupInfo 材料组请求结构(含id)
type MaterialGroupInfo struct {
	// MaterialGroupID string 材料组ID
	// 唯一标识材料组的UUID
	MaterialGroupID string `json:"materialGroupId" example:"123e4567-e89b-12d3-a456-426614174003"`

	// MaterialGroupParentID string 材料组父级ID 顶级材料组为空
	// 唯一标识材料组父级的UUID
	ParentID string `json:"parentId,omitempty" example:"123e4567-e89b-12d3-a456-426614174003"`

	// MaterialGroupName string 材料组名称
	// 材料组的名称信息
	MaterialGroupName string `json:"materialGroupName" example:"材料组名称"`

	// MaterialGroupCategory string 材料组类别
	// 材料组类别 resin=树脂, hardener=固化剂
	MaterialGroupCategory string `json:"materialGroupCategory" example:"材料组类别"`

	// Proportion float64 材料组占比
	// 材料组在实验步骤中的占比，百分比形式
	Proportion float64 `json:"proportion" example:"25.50"`

	// SubGroups []MaterialGroupInfo 子材料组
	// 当前材料组包含的子材料组列表，支持层级嵌套
	SubGroups []MaterialGroupInfo `json:"subGroups,omitempty"`

	// Materials []MaterialInfo 材料列表
	// 材料组内的具体材料信息
	Materials []MaterialInfo `json:"materials"`
}

// MaterialInfo 材料数据结构(含id)
type MaterialInfo struct {
	// MaterialID string 材料ID
	// 唯一标识材料的UUID
	MaterialID string `json:"materialId" example:"123e4567-e89b-12d3-a456-426614174004"`
	MaterialData
}

// MaterialData 材料数据结构
type MaterialData struct {
	// MaterialName string 材料名称
	// 材料的名称信息
	MaterialName string `json:"materialName" example:"材料名称"`

	// Percentage float64 材料占比
	// 材料在材料组中的占比，百分比形式
	Percentage float64 `json:"percentage" example:"60.00"`
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
	StartTime string `json:"startTime" validate:"omitempty,datetime=2006-01-02 15:04:05" example:"2024-02-05T08:30:00Z"`

	// EndTime string 实验结束时间
	// 选填，格式为 "2006-01-02T15:04:05Z"（RFC3339 格式）
	// 示例值: "2024-02-05T18:00:00Z"
	EndTime string `json:"endTime" validate:"omitempty,datetime=2006-01-02 15:04:05" example:"2024-02-05T18:00:00Z"`

	// RecipeID string 配方ID
	// 配方的UUID
	RecipeID string `json:"recipeId" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174002"`

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

	// StepCondition string 实验步骤的实验条件
	// 选填，实验步骤的实验条件描述，最长 255 字符
	// 示例值: "实验条件"
	StepCondition string `json:"stepCondition" validate:"omitempty,max=255" example:"实验步骤的实验条件"`

	// ResultValue string 步骤结果值
	// 选填，实验步骤的结果值，最长 256 字符
	// 示例值: "步骤结果值"
	ResultValue string `json:"resultValue" validate:"required,max=256" example:"步骤结果值"`

	// StepOrder int 排序(从大到小)
	// 必填，实验步骤的执行排序(从大到小)，必须为正整数
	// 示例值: 1
	StepOrder int `json:"stepOrder" validate:"required,gte=1" example:"1"`
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
	StartTime string `json:"startTime" validate:"omitempty,datetime=2006-01-02 15:04:05" example:"2024-02-05T08:30:00Z"`

	// EndTime string 实验结束时间
	// 选填，格式为 "2006-01-02T15:04:05Z"（RFC3339 格式）
	// 示例值: "2024-02-05T18:00:00Z"
	EndTime string `json:"endTime" validate:"omitempty,datetime=2006-01-02 15:04:05" example:"2024-02-05T18:00:00Z"`

	// RecipeID string 配方ID
	// 配方的UUID
	RecipeID string `json:"recipeId" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174002"`

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

	// StepCondition string 实验步骤的实验条件
	// 选填，实验步骤的实验条件描述，最长 255 字符
	// 示例值: "实验条件"
	StepCondition string `json:"stepCondition" validate:"omitempty,max=255" example:"实验步骤的实验条件"`

	// ResultValue string 步骤结果值
	// 选填，实验步骤的结果值，最长 256 字符
	// 示例值: "步骤结果值"
	ResultValue string `json:"resultValue" validate:"required,max=256" example:"步骤结果值"`

	// StepOrder int 排序(从大到小)
	// 必填，实验步骤的执行排序(从大到小)，必须为正整数
	// 示例值: 1
	StepOrder int `json:"stepOrder" validate:"required,gte=1" example:"1"`
}
