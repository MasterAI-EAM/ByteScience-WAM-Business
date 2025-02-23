package data

type PredictionRequest struct {
	// Steps []PredictionStepData 实验步骤
	// 包含该实验的步骤信息
	Steps []PredictionStepData `json:"steps"`
}

type PredictionStepData struct {
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

type PredictionResponse struct {
	// List []PredictionData 数据
	// 返回的实验记录结果
	List []PredictionData `json:"list"`
}

type PredictionData struct {
	// StepName string 步骤名称
	// 描述实验步骤的名称
	StepName string `json:"stepName" example:"步骤名称"`

	// ResultValue string 实验条件
	// 步骤对应的实验结果
	ResultValue string `json:"resultValue" example:"步骤结果值"`
}

type PredictionDemoRequest struct {
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

type PredictionDemoResponse struct {
	// StepName string 步骤名称
	// 描述实验步骤的名称
	StepName string `json:"stepName" example:"步骤名称"`

	// ResultValue string 实验条件
	// 步骤对应的实验结果
	ResultValue string `json:"resultValue" example:"步骤结果值"`
}
