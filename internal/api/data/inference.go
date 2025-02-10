package data

import (
	"ByteScience-WAM-Business/internal/model/dto/data"
	"ByteScience-WAM-Business/internal/service"
	"github.com/gin-gonic/gin"
)

type InferenceApi struct {
	service *service.InferenceService
}

// NewInferenceApi 创建 InferenceApi 实例并初始化依赖项
func NewInferenceApi() *InferenceApi {
	service := service.NewInferenceService()
	return &InferenceApi{service: service}
}

// Prediction 根据配方预测实验结果
// @Summary 根据配方预测实验结果
// @Description 通过提供的配方信息，预测实验结果
// @Tags 模型预测
// @Accept json
// @Produce json
// @Param body body data.PredictionRequest true "配方推荐请求参数"
// @Success 200 {object} data.PredictionResponse "推荐成功，返回推荐的材料列表"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，可能是配方信息不完整或格式不正确"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是预测服务异常"
// @Router /data/inference/prediction [post]
func (api InferenceApi) Prediction(ctx *gin.Context, req *data.PredictionRequest) (res *data.PredictionResponse,
	err error) {
	res, err = api.service.Prediction(ctx, req)
	return
}

// PredictionDemo 根据配方预测实验结果（Demo）
// @Summary 根据配方预测实验结果（Demo）
// @Description 演示版本的配方推荐功能，提供简单的预测功能
// @Tags 模型预测
// @Accept json
// @Produce json
// @Param body body data.PredictionDemoRequest true "配方推荐演示请求参数"
// @Success 200 {object} data.PredictionDemoResponse "演示推荐成功，返回推荐的材料示例"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，可能是配方信息不完整或格式不正确"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是演示服务异常"
// @Router /data/inference/prediction/demo [post]
func (api InferenceApi) PredictionDemo(ctx *gin.Context, req *data.PredictionDemoRequest) (res *data.
	PredictionDemoResponse,
	err error) {
	res, err = api.service.PredictionDemo(ctx, req)
	return
}
