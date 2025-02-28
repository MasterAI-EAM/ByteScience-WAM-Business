package ai

import (
	"ByteScience-WAM-Business/internal/model/dto/ai"
	"ByteScience-WAM-Business/internal/service"
	"ByteScience-WAM-Business/internal/utils"
	"github.com/gin-gonic/gin"
)

type PredictApi struct {
	service *service.PredictService
}

// NewPredictApi 创建 PredictApi 实例并初始化依赖项
func NewPredictApi() *PredictApi {
	service := service.NewPredictService()
	return &PredictApi{service: service}
}

// ForwardDirection 根据配方信息预测实验结果
// @Summary 根据配方信息预测实验结果
// @Description 提供配方信息并进行预测，返回实验结果的预测值
// @Tags 模型预测
// @Accept json
// @Produce json
// @Param body body ai.ForwardDirectionRequest true "配方预测请求参数"
// @Success 200 {object} ai.ForwardDirectionResponse "预测成功，返回预测的实验结果"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，可能是配方信息不完整或格式不正确"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是预测服务异常"
// @Router /ai/predict/forwardDirection [post]
func (api PredictApi) ForwardDirection(ctx *gin.Context, req *ai.ForwardDirectionRequest) (res *ai.ForwardDirectionResponse,
	err error) {
	// 百分比校验
	var totalProportion, totalPercentage float64
	for _, group := range req.MaterialGroups {
		totalProportion += group.Proportion
		for _, material := range group.Materials {
			totalPercentage += material.Percentage
		}
		if totalPercentage != 100 {
			return nil, utils.NewBusinessError(utils.MaterialProportionSumNot100Code, "")
		}
		totalPercentage = 0
	}
	if totalProportion != 100 {
		return nil, utils.NewBusinessError(utils.MaterialGroupProportionNot100Code, "")
	}

	res, err = api.service.ForwardDirection(ctx, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return
}
