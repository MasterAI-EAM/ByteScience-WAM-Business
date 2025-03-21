package service

import (
	"ByteScience-WAM-Business/internal/dao"
	"ByteScience-WAM-Business/internal/model/dto/ai"
	"ByteScience-WAM-Business/internal/utils"
	"context"
)

type PredictService struct {
	experimentDao *dao.ExperimentDao
}

// NewPredictService 创建一个新的 PredictService 实例
func NewPredictService() *PredictService {
	return &PredictService{
		experimentDao: dao.NewExperimentDao(),
	}
}

// ForwardDirection 根据配方预测结果
func (ps *PredictService) ForwardDirection(ctx context.Context,
	req *ai.ForwardDirectionRequest) (*ai.ForwardDirectionResponse, error) {

	res := &ai.ForwardDirectionResponse{
		HistoryList: make([]ai.ForwardDirectionResultInfo, 0),
	}
	//
	// // 格式化数据
	// var materials []entity.Materials
	// var materialGroups []entity.MaterialGroups
	// recipeID := uuid.NewString()
	// for _, materialGroup := range req.MaterialGroups {
	// 	materialGroupData := entity.MaterialGroups{
	// 		RecipeID:          recipeID,
	// 		ID:                uuid.NewString(),
	// 		Proportion:        materialGroup.Proportion,
	// 		MaterialGroupName: materialGroup.MaterialGroupName,
	// 	}
	// 	materialGroups = append(materialGroups, materialGroupData)
	// 	for _, row := range materialGroup.Materials {
	// 		material := entity.Materials{
	// 			ID:              uuid.NewString(),
	// 			MaterialName:    row.MaterialName,
	// 			MaterialGroupID: materialGroupData.ID,
	// 			Proportion:      row.Proportion,
	// 		}
	// 		materials = append(materials, material)
	// 	}
	// }
	//
	// // 生成配方密钥
	// recipeSignature := GenerateRecipeSignature(req.ExperimentCondition, materialGroups, materials)
	//
	// // 检索数据
	// result, err := ps.recipeDao.GetDataByRecipeSignature(ctx, recipeSignature, req.StepName, req.ExperimentCondition)
	// if err != nil {
	// 	logger.Logger.Errorf("[PredictService ForwardDirection] GetRecipeInSignatureMap err: %v", err)
	// 	return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	// }
	//
	// if result != nil && len(result) > 0 {
	// 	experimentIdList := make([]string, 0)
	// 	for _, step := range result {
	// 		experimentIdList = append(experimentIdList, step.ExperimentID)
	// 	}
	//
	// 	experimentList, err := ps.experimentDao.GetByIDList(ctx, experimentIdList)
	// 	if err != nil {
	// 		logger.Logger.Errorf("[PredictService ForwardDirection] GetRecipeInSignatureMap err: %v", err)
	// 		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	// 	}
	//
	// 	experimentMap := make(map[string]string)
	// 	for _, experiment := range experimentList {
	// 		experimentMap[experiment.ID] = experiment.ExperimentName
	// 	}
	// 	for _, step := range result {
	// 		res.HistoryList = append(res.HistoryList, ai.ForwardDirectionResultInfo{
	// 			ExperimentName: experimentMap[step.ExperimentID],
	// 			ForwardDirectionResult: ai.ForwardDirectionResult{
	// 				StepName:    step.StepName,
	// 				ResultValue: step.ResultValue,
	// 			},
	// 		})
	// 	}
	// 	return res, nil
	// }

	var err error
	// 调用ai接口查询
	res.AiResult, err = utils.SendPredictionRequest(req)
	if err != nil {
		return nil, err
	}

	// 返回最终结果
	return res, nil
}
