package ai

import "ByteScience-WAM-Business/internal/model/dto/data"

type ForwardDirectionRequest struct {
	// StepName string 步骤名称
	// 描述实验步骤的名称
	StepName string `json:"stepName" example:"步骤名称"`

	// ExperimentCondition string 实验条件
	// 步骤对应的实验条件描述
	ExperimentCondition string `json:"experimentCondition" example:"实验条件"`

	// MaterialGroups []MaterialGroupData 材料组
	// 该步骤中涉及的材料组信息
	MaterialGroups []data.MaterialGroupData `json:"materialGroups"`
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
