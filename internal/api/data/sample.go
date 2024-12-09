package data

import (
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/model/dto/data"
	"ByteScience-WAM-Business/internal/service"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/logger"
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

// Import 导入文件
// @Summary 导入文件
// @Description 接收上传的文件并处理，根据业务需求进行相关文件解析和导入
// @Tags 实验管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Success 200 {object} dto.Empty "文件导入成功，返回空对象表示操作成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，可能是文件上传失败或格式不符合要求"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是文件解析或存储过程中出现异常"
// @Router /data/sample/import [post]
func (api *SampleApi) Import(ctx *gin.Context, req *dto.Empty) (res *dto.Empty, err error) {
	// 获取文件
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Logger.Errorf("[Import] ctx.FormFile Error: %v", err)
		return nil, utils.NewBusinessError(utils.FileParsingFailedCode)
	}

	// 文件不能超出10m
	if file.Size > 10*1024*1024 {
		logger.Logger.Warnf("File size exceeds 10 MB: %d bytes", file.Size)
		return nil, utils.NewBusinessError(utils.FileTooLargeCode)
	}

	res, err = api.service.Import(ctx, file)
	if err != nil {
		logger.Logger.Errorf("[Import] Import Error: %v", err)
		return nil, utils.NewBusinessError(utils.FileParsingFailedCode)
	}

	return
}

// List 获取实验列表
// @Summary 获取实验列表
// @Description 根据分页请求获取实验列表，支持按实验名称进行筛选
// @Tags 实验管理
// @Accept json
// @Produce json
// @Param req body data.ExperimentListRequest true "请求参数，包含分页信息及筛选条件"
// @Success 200 {object} data.ExperimentListResponse "成功获取实验列表，返回实验数据"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如分页参数错误、筛选条件不符合要求"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库查询出错等情况"
// @Router /data/sample [get]
func (api *SampleApi) List(ctx *gin.Context, req *data.ExperimentListRequest) (res *data.ExperimentListResponse, err error) {
	res, err = api.service.List(ctx, req)
	if err != nil {
		logger.Logger.Errorf("[List] List Error: %v", err)
		return nil, utils.NewBusinessError(utils.FileParsingFailedCode)
	}
	return
}

// Delete 删除实验
// @Summary 删除实验
// @Description 根据实验ID删除实验
// @Tags 实验管理
// @Accept json
// @Produce json
// @Param req body data.ExperimentDeleteRequest true "请求参数，包含要删除的实验ID"
// @Success 200 {object} dto.Empty "成功删除实验"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如实验ID不存在或格式无效"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库删除失败等情况"
// @Router /data/sample [delete]
func (api *SampleApi) Delete(ctx *gin.Context, req *data.ExperimentDeleteRequest) (res *dto.Empty, err error) {
	res, err = api.service.Delete(ctx, req)
	return
}

// Edit 修改实验信息
// @Summary 修改实验信息
// @Description 根据实验ID修改实验的具体信息
// @Tags 实验管理
// @Accept json
// @Produce json
// @Param req body data.ExperimentUpdateRequest true "请求参数，包含实验ID及要修改的具体内容"
// @Success 200 {object} dto.Empty "成功修改实验信息"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如实验ID不存在或修改内容无效"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库更新失败等情况"
// @Router /data/sample [put]
func (api *SampleApi) Edit(ctx *gin.Context, req *data.ExperimentUpdateRequest) (res *dto.Empty, err error) {
	res, err = api.service.Edit(ctx, req)
	return
}
