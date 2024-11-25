package service

import (
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/model/dto/auth"
	"context"
)

type SampleService struct {
}

// NewSampleService 创建一个新的 SampleService 实例
func NewSampleService() *SampleService {
	return &SampleService{}
}

// List 登录方法
func (ds SampleService) List(ctx context.Context, req *auth.LoginRequest) (*dto.Empty, error) {
	return nil, nil
}
