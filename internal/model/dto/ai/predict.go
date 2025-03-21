package ai

type ForwardDirectionRequest struct {
	// StepName string 步骤名称
	// 描述实验步骤的名称
	StepName string `json:"stepName" example:"步骤名称"`

	// ExperimentCondition string 实验条件
	// 步骤对应的实验条件描述
	ExperimentCondition string `json:"experimentCondition" example:"实验条件"`

	// MaterialGroups []MaterialGroupData 材料组
	// 该步骤中涉及的材料组信息
	MaterialGroups []MaterialGroupData `json:"materialGroups"`
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

	// Proportion float64 材料占比
	// 材料在材料组中的占比，百分比形式
	Proportion float64 `json:"proportion" example:"25.50"`
}

type ForwardDirectionResponse struct {
	// AiResult ForwardDirectionResult 数据
	// ai的结果
	AiResult *ForwardDirectionResult `json:"aiResult"`

	// historyList []ForwardDirectionResultInfo 数据
	// 实验记录的结果
	HistoryList []ForwardDirectionResultInfo `json:"historyList"`
}

type ForwardDirectionResultInfo struct {
	// ExperimentName string 实验名称
	// 实验的名称信息
	ExperimentName string `json:"experimentName" example:"实验名称"`
	ForwardDirectionResult
}

type ForwardDirectionResult struct {
	// StepName string 步骤名称
	// 描述实验步骤的名称
	StepName string `json:"stepName" example:"步骤名称"`

	// ResultValue string 实验条件
	// 步骤对应的实验结果
	ResultValue string `json:"resultValue" example:"步骤结果值"`
}
