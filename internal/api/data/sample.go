package data

import (
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/service"
	"github.com/gin-gonic/gin"
)

type SampleApi struct {
	service *service.SampleService
}

// NewSampleApi 创建 SampleApi 实例并初始化依赖项
func NewSampleApi() *SampleApi {
	service := service.NewSampleService()
	return &SampleApi{service: service}
}

func (api *SampleApi) List(ctx *gin.Context, req *dto.Empty) (res *dto.Empty, err error) {
	return
}
