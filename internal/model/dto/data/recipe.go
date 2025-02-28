package data

import "ByteScience-WAM-Business/internal/model/dto"

// RecipeListRequest 用于分页查询配方列表的请求体结构
type RecipeListRequest struct {
	// PaginationRequest dto.PaginationRequest 通用分页参数
	// 包含页码和每页大小等分页信息
	dto.PaginationRequest

	// RecipeName string 配方名称，选填，长度限制：2-255字符
	// 用于按名称查询配方列表 前缀模糊匹配
	RecipeName string `json:"recipeName" validate:"omitempty,min=2,max=255" example:"配方名称"`
}

// RecipeListResponse 分页查询配方列表的响应体结构
type RecipeListResponse struct {
	// Total int64 总条数
	// 返回符合条件的实验记录总数
	Total int64 `json:"total" example:"100"`

	// List []RecipeData 数据
	// 分页返回的配方列表
	List []RecipeData `json:"list"`
}

type RecipeData struct {
	// RecipeId string 配方id
	// UUID
	RecipeId string `json:"recipeId" example:"id"`

	// RecipeName string 配方名称
	// 配方名称信息
	RecipeName string `json:"recipeName" example:"配方名称"`

	// RecipeUsedInExperimentNum int64
	// 配方被实验使用数
	RecipeUsedInExperimentNum int64 `json:"recipeUsedInExperimentNum" example:"11"`

	// Sort int
	// 排序 优先级从大到小
	Sort int `json:"sort" example:"1"`

	// IsErr bool 是否发生错误
	// 表示该配方数据处理过程中是否出现错误，true 表示有错误，false 表示无错误
	IsErr bool `json:"isErr" example:"false"`

	// ErrMsg string 错误信息
	// 当 isErr 为 true 时，该字段包含具体的错误描述信息
	ErrMsg string `json:"errMsg" example:"The proportion of the material group is not 100%"`

	// CreatedAt 创建时间
	// 格式为时间戳，创建时间
	CreatedAt string `json:"createdAt" example:"2024-11-18T10:00:00Z"`

	// MaterialGroups []MaterialGroupInfo 材料组
	// 该步骤中涉及的材料组信息
	MaterialGroups []MaterialGroupInfo `json:"materialGroups"`
}

// RecipeFormListRequest 表示用于获取配方列表的请求结构体。
type RecipeFormListRequest struct {
	// PaginationRequest 是通用的分页请求参数结构体，嵌套在本请求结构体中。
	// 它包含页码（Page）和每页大小（PageSize）等分页信息，用于控制返回结果的分页情况。
	// 例如，当 Page 为 1，PageSize 为 10 时，将返回第一页的 10 条记录。
	dto.PaginationRequest

	// RecipeName 是一个可选的字符串字段，用于按配方名称进行模糊查询。
	// 该字段的长度必须在 2 到 255 个字符之间。如果不提供该字段，则不进行名称过滤。
	// 示例值："配方名称"，系统会查找名称中包含该字符串的所有配方记录。
	RecipeName string `json:"recipeName" validate:"omitempty,min=2,max=255" example:"配方名称"`
}

// RecipeFormListResponse 表示获取配方列表请求的响应结构体。
type RecipeFormListResponse struct {
	// Total 表示符合查询条件的配方记录的总条数。
	// 示例值："100"，表示共有 100 条记录符合查询条件。
	Total int64 `json:"total" example:"100"`

	// List 是一个 RecipeInfo 类型的切片，用于存储分页返回的配方记录列表。
	// 例如，列表中可能包含多个不同配方的信息，前端可以根据这些信息展示配方列表。
	List []RecipeInfo `json:"list"`
}

// RecipeInfo 表示单个配方的详细信息结构体。
type RecipeInfo struct {
	// Id 是配方的唯一标识符，通常为 UUID 格式的字符串。
	// 示例值: "123e4567-e89b-12d3-a456-426614174000"
	Id string `json:"id"`

	// Name 是配方的名称，用于描述该配方的内容或用途。
	// 示例值: "实验名称"
	Name string `json:"name"`
}

// RecipeInfoRequest 用于查询配方详情的请求体结构
type RecipeInfoRequest struct {
	// RecipeId 是配方的唯一标识符
	// 必填字段，格式为 UUID v4
	// 示例值: "123e4567-e89b-12d3-a456-426614174000"
	RecipeId string `json:"recipeId" validate:"required,uuid4" example:"配方id"`
}

// RecipeInfoResponse 表示查询配方详情的响应体结构
type RecipeInfoResponse struct {
	// RecipeId string 配方id
	// UUID
	RecipeId string `json:"recipeId" example:"id"`

	// RecipeName 表示配方的名称
	// 返回配方的名称信息
	// 示例值: "配方名称"
	RecipeName string `json:"recipeName" example:"配方名称"`

	// CreatedAt 表示配方的创建时间
	// 时间格式为 ISO 8601 格式的 UTC 时间戳
	// 示例值: "2024-11-18T10:00:00Z"
	CreatedAt string `json:"createdAt" example:"2024-11-18T10:00:00Z"`

	// MaterialGroups 包含与该配方关联的材料组信息
	// 每个材料组都包含详细的材料信息
	MaterialGroups []MaterialGroupInfo `json:"materialGroups"`

	// RecipeBasedExperiment 包含基于此配方创建的实验信息列表
	RecipeBasedExperiment []ExperimentInfo `json:"recipeBasedExperiment"`
}

// ExperimentInfo 表示与配方关联的实验信息
type ExperimentInfo struct {
	// Id 是实验的唯一标识符
	// 示例值: "123e4567-e89b-12d3-a456-426614174000"
	Id string `json:"id"`

	// Name 是实验的名称
	// 示例值: "实验名称"
	Name string `json:"name"`
}

// MaterialGroupData 材料组数据结构
type MaterialGroupData struct {
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
	// MaterialName string 材料名称
	// 材料的名称信息
	MaterialName string `json:"materialName" example:"材料名称"`

	// Percentage float64 材料占比
	// 材料在材料组中的占比，百分比形式
	Percentage float64 `json:"percentage" example:"60.00"`
}

// MaterialGroupInfo 材料组请求结构(含id)
type MaterialGroupInfo struct {
	// MaterialGroupID string 材料组ID
	// 唯一标识材料组的UUID
	MaterialGroupID string `json:"materialGroupId" example:"123e4567-e89b-12d3-a456-426614174003"`

	// MaterialGroupName string 材料组名称
	// 材料组的名称信息
	MaterialGroupName string `json:"materialGroupName" example:"材料组名称"`

	// Proportion float64 材料组占比
	// 材料组在实验步骤中的占比，百分比形式
	Proportion float64 `json:"proportion" example:"25.50"`

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

// RecipeAddRequest 添加配方的请求体结构
type RecipeAddRequest struct {
	// RecipeName string 配方名称，选填，长度限制：2-255字符
	// 用于按名称模糊查询实验记录
	RecipeName string `json:"recipeName" validate:"omitempty,min=2,max=255" example:"配方名称"`

	// Sort int 排序 优先级从大到小
	// 选填，必须是大于等于 0 的整数
	// 示例值: 1
	Sort int `json:"sort" validate:"omitempty,min=0" example:"1"`

	// MaterialGroupData []MaterialGroupData 材料组列表
	// 选填，实验步骤中涉及的材料组信息, percentage的和为100(占比100%)
	// 示例值: [{"materialGroupName": "材料组名称", "proportion": 25.5, "materials": []}]
	MaterialGroups []MaterialGroupData `json:"materialGroups" validate:"omitempty,dive"`
}

// RecipeDeleteRequest 用于删除配方详情的请求体结构
type RecipeDeleteRequest struct {
	// RecipeId 是配方的唯一标识符
	// 必填字段，格式为 UUID v4
	// 示例值: "123e4567-e89b-12d3-a456-426614174000"
	RecipeId string `json:"recipeId" validate:"required,uuid4" example:"配方id"`
}

// RecipeEditRequest 用于编辑配方详情的请求体结构
type RecipeEditRequest struct {
	// RecipeId 是配方的唯一标识符
	// 必填字段，格式为 UUID v4
	// 示例值: "123e4567-e89b-12d3-a456-426614174000"
	RecipeId string `json:"recipeId" validate:"required,uuid4" example:"配方id"`

	// RecipeName string 配方名称，选填，长度限制：2-255字符
	// 用于按名称模糊查询实验记录
	RecipeName string `json:"recipeName" validate:"required,min=2,max=255" example:"配方名称"`

	// Sort int 排序 优先级从大到小
	// 选填，必须是大于等于 0 的整数
	// 示例值: 1
	Sort int `json:"sort" validate:"required,min=0" example:"1"`

	// MaterialGroupData []MaterialGroupData 材料组列表
	// 选填，实验步骤中涉及的材料组信息, percentage的和为100(占比100%)
	// 示例值: [{"materialGroupName": "材料组名称", "proportion": 25.5, "materials": []}]
	MaterialGroups []MaterialGroupInfo `json:"materialGroups" validate:"required,dive"`
}
