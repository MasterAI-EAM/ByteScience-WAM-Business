package data

import (
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/model/dto/data"
	"ByteScience-WAM-Business/internal/service"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ExperimentApi struct {
	service *service.ExperimentService
}

// NewExperimentApi 创建 ExperimentApi 实例并初始化依赖项
func NewExperimentApi() *ExperimentApi {
	service := service.NewExperimentService()
	return &ExperimentApi{service: service}
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
// @Router /data/experiment [get]
func (api *ExperimentApi) List(ctx *gin.Context, req *data.ExperimentListRequest) (res *data.ExperimentListResponse, err error) {
	res, err = api.service.List(ctx, req)
	if err != nil {
		logger.Logger.Errorf("[List] List Error: %v", err)
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
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
// @Router /data/experiment [delete]
func (api *ExperimentApi) Delete(ctx *gin.Context, req *data.ExperimentDeleteRequest) (res *dto.Empty, err error) {
	userId, exists := ctx.Get("userId")
	if !exists {
		logger.Logger.Errorf("[Import] User id does not exist")
		return nil, utils.NewBusinessError(utils.UserNotFoundCode, "")
	}

	res, err = api.service.Delete(ctx, userId.(string), req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return
}

// Info 获取实验详情
// @Summary 获取指定实验的详情信息
// @Description 根据实验ID获取实验的详细信息
// @Tags 实验管理
// @Accept json
// @Produce json
// @Param req body data.ExperimentInfoRequest true "请求参数，包含实验ID"
// @Success 200 {object} data.ExperimentInfoResponse "成功返回实验详情"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如实验ID不存在或格式无效"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是查询失败等情况"
// @Router /data/experiment/info [get]
func (api *ExperimentApi) Info(ctx *gin.Context, req *data.ExperimentInfoRequest) (res *data.ExperimentInfoResponse, err error) {
	res, err = api.service.Info(ctx, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return
}

//
// // Add 添加实验信息
// // @Summary 添加新的实验信息
// // @Description 根据提供的实验信息创建新的实验记录
// // @Tags 实验管理
// // @Accept json
// // @Produce json
// // @Param req body data.ExperimentAddRequest true "请求参数，包含新的实验信息"
// // @Success 200 {object} dto.Empty "成功添加实验信息"
// // @Failure 400 {object} dto.ErrorResponse "请求参数错误，如缺少必要字段或格式不正确"
// // @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库插入失败等情况"
// // @Router /data/experiment [post]
// func (api *ExperimentApi) Add(ctx *gin.Context, req *data.ExperimentAddRequest) (res *dto.Empty, err error) {
// 	userId, exists := ctx.Get("userId")
// 	if !exists {
// 		logger.Logger.Errorf("[Import] User id does not exist")
// 		return nil, utils.NewBusinessError(utils.UserNotFoundCode, "")
// 	}
//
// 	// 调用服务层的 Add 方法进行实验数据插入
// 	res, err = api.service.Add(ctx, userId.(string), req)
// 	if err != nil {
// 		if businessErr, ok := err.(*utils.BusinessError); ok {
// 			return nil, businessErr
// 		}
// 		return nil, utils.NewBusinessError(utils.InternalError, "")
// 	}
//
// 	return res, err
// }
//
// // Edit 修改实验信息
// // @Summary 修改实验信息
// // @Description 根据实验ID修改实验的具体信息
// // @Tags 实验管理
// // @Accept json
// // @Produce json
// // @Param req body data.ExperimentUpdateRequest true "请求参数，包含实验ID及要修改的具体内容"
// // @Success 200 {object} dto.Empty "成功修改实验信息"
// // @Failure 400 {object} dto.ErrorResponse "请求参数错误，如实验ID不存在或修改内容无效"
// // @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库更新失败等情况"
// // @Router /data/experiment [put]
// func (api *ExperimentApi) Edit(ctx *gin.Context, req *data.ExperimentUpdateRequest) (res *dto.Empty, err error) {
// 	userId, exists := ctx.Get("userId")
// 	if !exists {
// 		logger.Logger.Errorf("[Import] User id does not exist")
// 		return nil, utils.NewBusinessError(utils.UserNotFoundCode, "")
// 	}
//
// 	res, err = api.service.Edit(ctx, userId.(string), req)
// 	if err != nil {
// 		if businessErr, ok := err.(*utils.BusinessError); ok {
// 			return nil, businessErr
// 		}
// 		return nil, utils.NewBusinessError(utils.InternalError, "")
// 	}
//
// 	return
// }
