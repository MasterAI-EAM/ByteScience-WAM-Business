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

	// Steps []ExperimentStepData 实验步骤
	// 包含该实验的步骤信息
	Steps []ExperimentStepData `json:"steps"`
}

// ExperimentStepData 实验步骤数据结构
type ExperimentStepData struct {
	// StepID string 步骤ID
	// 唯一标识实验步骤的UUID
	StepID string `json:"stepId" example:"123e4567-e89b-12d3-a456-426614174002"`

	// StepName string 步骤名称
	// 描述实验步骤的名称
	StepName string `json:"stepName" example:"步骤名称"`

	// StepNameDescription string 实验步骤描述
	// 实验步骤描述
	StepNameDescription string `json:"stepNameDescription" example:"实验步骤描述"`

	// ExperimentCondition string 实验条件
	// 步骤对应的实验条件描述
	ExperimentCondition string `json:"experimentCondition" example:"实验条件"`

	// ResultValue string 实验条件
	// 步骤对应的实验结果
	ResultValue string `json:"resultValue" example:"步骤结果值"`

	// MaterialGroups []MaterialGroupData 材料组
	// 该步骤中涉及的材料组信息
	MaterialGroups []MaterialGroupData `json:"materialGroups"`
}

// MaterialGroupData 材料组数据结构
type MaterialGroupData struct {
	// MaterialGroupID string 材料组ID
	// 唯一标识材料组的UUID
	MaterialGroupID string `json:"materialGroupId" example:"123e4567-e89b-12d3-a456-426614174003"`

	// MaterialGroupName string 材料组名称
	// 材料组的名称信息
	MaterialGroupName string `json:"materialGroupName" example:"材料组名称"`

	// Proportion float64 材料组占比
	// 材料组在实验步骤中的占比，百分比形式
	Proportion float64 `json:"proportion" example:"25.50"`

	// Materials []MaterialData 材料列表
	// 材料组内的具体材料信息
	Materials []MaterialData `json:"materials"`
}

// MaterialData 材料数据结构
type MaterialData struct {
	// MaterialID string 材料ID
	// 唯一标识材料的UUID
	MaterialID string `json:"materialId" example:"123e4567-e89b-12d3-a456-426614174004"`

	// MaterialName string 材料名称
	// 材料的名称信息
	MaterialName string `json:"materialName" example:"材料名称"`

	// Percentage float64 材料占比
	// 材料在材料组中的占比，百分比形式
	Percentage float64 `json:"percentage" example:"60.00"`
}
