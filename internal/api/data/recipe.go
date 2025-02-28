package data

import (
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/model/dto/data"
	"ByteScience-WAM-Business/internal/service"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/logger"
	"github.com/gin-gonic/gin"
)

type RecipeApi struct {
	service *service.RecipeService
}

// NewRecipeApi 创建 RecipeApi 实例并初始化依赖项
func NewRecipeApi() *RecipeApi {
	service := service.NewRecipeService()
	return &RecipeApi{service: service}
}

// List 获取配方列表
// @Summary 获取配方列表
// @Description 根据分页请求获取配方列表，支持按配方名称进行筛选
// @Tags 配方管理
// @Accept json
// @Produce json
// @Param req body data.RecipeListRequest true "请求参数，包含分页信息及筛选条件"
// @Success 200 {object} data.RecipeListResponse "成功获取配方列表，返回配方数据"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如分页参数错误、筛选条件不符合要求"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库查询失败等情况"
// @Router /data/recipe [get]
func (api *RecipeApi) List(ctx *gin.Context, req *data.RecipeListRequest) (res *data.RecipeListResponse, err error) {
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

// Add 添加新的配方信息
// @Summary 添加新的配方
// @Description 根据提供的配方信息创建新的配方记录
// @Tags 配方管理
// @Accept json
// @Produce json
// @Param req body data.RecipeAddRequest true "请求参数，包含新的配方信息"
// @Success 200 {object} dto.Empty "成功添加配方信息"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如缺少必要字段或格式不正确"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库插入失败等情况"
// @Router /data/recipe [post]
func (api *RecipeApi) Add(ctx *gin.Context, req *data.RecipeAddRequest) (res *dto.Empty, err error) {
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
	// 调用服务层的 Add 方法进行配方数据插入
	res, err = api.service.Add(ctx, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return res, err
}

// Info 获取单个配方详情
// @Summary 获取单个配方详情
// @Description 根据传入的配方 ID 获取该配方的详细信息，包括材料组、材料信息以及基于此配方创建的实验信息
// @Tags 配方管理
// @Accept json
// @Produce json
// @Param req body data.RecipeInfoRequest true "请求参数，包含要查询的配方 ID"
// @Success 200 {object} data.RecipeInfoResponse "成功获取单个配方详情，返回配方详细数据"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如配方 ID 格式错误等"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库查询失败等情况"
// @Router /data/recipe/info [get]
func (api *RecipeApi) Info(ctx *gin.Context, req *data.RecipeInfoRequest) (res *data.RecipeInfoResponse, err error) {
	res, err = api.service.Info(ctx, req)
	if err != nil {
		logger.Logger.Errorf("[Info] Info Error: %v", err)
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return
}

// FormList 获取配方表单列表
// @Summary 获取配方表单列表
// @Description 根据分页请求获取配方表单列表，支持按配方名称进行筛选
// @Tags 配方管理
// @Accept json
// @Produce json
// @Param req body data.RecipeFormListRequest true "请求参数，包含分页信息及筛选条件"
// @Success 200 {object} data.RecipeFormListResponse "成功获取配方表单列表，返回配方数据"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如分页参数错误、筛选条件不符合要求"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库查询失败等情况"
// @Router /data/recipe/form/list [get]
func (api *RecipeApi) FormList(ctx *gin.Context, req *data.RecipeFormListRequest) (res *data.RecipeFormListResponse, err error) {
	res, err = api.service.FormList(ctx, req)
	if err != nil {
		logger.Logger.Errorf("[FormList] FormList Error: %v", err)
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return
}

// Edit 编辑现有配方信息
// @Summary 编辑现有配方
// @Description 根据提供的配方信息更新现有配方记录，更新时确保材料组和材料百分比的比例和为100%
// @Tags 配方管理
// @Accept json
// @Produce json
// @Param req body data.RecipeEditRequest true "请求参数，包含更新的配方信息"
// @Success 200 {object} dto.Empty "成功更新配方信息"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如缺少必要字段或格式不正确"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库更新失败等情况"
// @Router /data/recipe [put]
func (api *RecipeApi) Edit(ctx *gin.Context, req *data.RecipeEditRequest) (res *dto.Empty, err error) {
	// 校验材料组和材料百分比的和是否为100%
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

	res, err = api.service.Edit(ctx, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return res, err
}

// Delete 删除配方信息
// @Summary 删除配方
// @Description 根据配方ID删除指定的配方记录
// @Tags 配方管理
// @Accept json
// @Produce json
// @Param req body data.RecipeDeleteRequest true "请求参数，包含需要删除的配方ID"
// @Success 200 {object} dto.Empty "成功删除配方信息"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，可能是配方ID无效"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库删除失败等情况"
// @Router /data/recipe [delete]
func (api *RecipeApi) Delete(ctx *gin.Context, req *data.RecipeDeleteRequest) (res *dto.Empty, err error) {
	res, err = api.service.Delete(ctx, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return res, err
}
