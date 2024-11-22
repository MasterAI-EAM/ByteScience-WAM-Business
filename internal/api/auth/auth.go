package auth

import (
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/model/dto/auth"
	"ByteScience-WAM-Business/internal/service"
	"github.com/gin-gonic/gin"
)

type Api struct {
	service *service.AuthService
}

// NewAuthApi 创建 adminApi 实例并初始化依赖项
func NewAuthApi() *Api {
	service := service.NewAuthService()
	return &Api{service: service}
}

// Login 用户登录
// @Summary 用户登录
// @Description 提供用户名和密码进行登录操作，验证用户身份并获取相应权限
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param req body auth.LoginRequest true "请求参数，包含用户名、密码等登录所需信息"
// @Success 200 {object} auth.LoginResponse "成功登录，返回token凭证表示操作成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如用户名或密码格式不正确等"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库查询出错、验证逻辑异常等情况"
// @Router /login [post]
func (api *Api) Login(ctx *gin.Context, req *auth.LoginRequest) (res *auth.LoginResponse, err error) {
	res, err = api.service.Login(ctx, req)
	return
}

// ChangPassword 修改密码
// @Summary 修改用户密码
// @Description 根据提供的原密码及新密码等信息修改用户当前账户的密码
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param req body auth.ChangePasswordRequest true "请求参数，包含原密码、新密码等修改密码所需信息"
// @Success 200 {object} dto.Empty "成功修改密码，返回空对象表示操作成功"
// @Success 201 {object} dto.Empty "可根据实际情况设置不同成功状态码及对应含义，这里示例201可表示密码已成功更新并生效"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如原密码错误、新密码格式不符合要求等"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库更新出错、验证逻辑异常等情况"
// @Router /changPassword [put]
func (api *Api) ChangPassword(ctx *gin.Context, req *auth.ChangePasswordRequest) (res *dto.Empty, err error) {
	err = api.service.ChangePassword(ctx, req)
	return
}
