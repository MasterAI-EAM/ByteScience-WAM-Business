package service

import (
	"ByteScience-WAM-Business/internal/model/dto/data"
	"context"
)

type InferenceService struct{}

// NewInferenceService 创建一个新的 InferenceService 实例
func NewInferenceService() *InferenceService {
	return &InferenceService{}
}

// Prediction 根据配方推荐材料
func (is *InferenceService) Prediction(ctx context.Context, req *data.PredictionRequest) (*data.PredictionResponse, error) {

	if req == nil || len(req.Steps) == 0 {
		return nil, nil
	}

	list := make([]data.PredictionData, len(req.Steps))
	for k, v := range req.Steps {
		list[k] = data.PredictionData{
			StepName:    v.StepName,
			ResultValue: "123",
		}
	}

	// 返回最终结果
	return &data.PredictionResponse{List: list}, nil
}

// PredictionDemo 根据配方推荐材料
func (is *InferenceService) PredictionDemo(ctx context.Context, req *data.PredictionDemoRequest) (*data.PredictionDemoResponse, error) {

	// 返回最终结果
	return &data.PredictionDemoResponse{}, nil
}
