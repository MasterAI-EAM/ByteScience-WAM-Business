package service

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/internal/dao"
	"ByteScience-WAM-Business/internal/model/dto/auth"
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/logger"
	"context"
	"time"
)

type AuthService struct {
	dao *dao.UserDao // 添加 AdminDao 作为成员
}

// NewAuthService 创建一个新的 AuthService 实例
func NewAuthService() *AuthService {
	return &AuthService{
		dao: dao.NewUserDao(),
	}
}

// Login 登录方法
func (as AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// 确定 Identifier 类型（用户名、邮箱或手机号）
	identifierType := utils.IdentifyType(req.Identifier)
	var admin *entity.Users
	var err error

	switch identifierType {
	case "email":
		admin, err = as.dao.GetByFields(ctx, "", req.Identifier, "")
	case "phone":
		admin, err = as.dao.GetByFields(ctx, "", "", req.Identifier)
	default:
		admin, err = as.dao.GetByFields(ctx, req.Identifier, "", "")
	}

	// 检查用户是否存在
	if err != nil {
		logger.Logger.Errorf("[Login] Error fetching user by %s: %v", identifierType, err)
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	if admin == nil {
		return nil, utils.NewBusinessError(utils.UserNotFoundCode, "")
	}

	if admin.Status == 0 {
		return nil, utils.NewBusinessError(utils.UserNotFoundCode, "")
	}

	// 验证密码是否正确
	isMatch, err := utils.VerifyPassword(req.Password, admin.Password)
	if err != nil {
		// 如果发生了错误（非匹配错误），记录日志并返回
		logger.Logger.Errorf("[Login] Error verifying password: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	if !isMatch {
		// 如果密码不匹配，返回无效凭证错误
		return nil, utils.NewBusinessError(utils.PasswordIncorrectCode, "")
	}

	// 生成 JWT Token
	token, err := utils.GetToken(conf.GlobalConf.Jwt.AccessSecret, conf.GlobalConf.Jwt.AccessExpire, admin.ID)
	if err != nil {
		logger.Logger.Errorf("[Login] Error utils.GetToken: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}

	// 记录登陆时间
	if err = as.dao.UpdateLastLoginTime(ctx, admin.ID); err != nil {
		logger.Logger.Errorf("[Login] Error UpdateLastLoginTime: %v", err)
	}

	// 返回登录响应
	loginResponse := &auth.LoginResponse{
		Token: token,
	}

	return loginResponse, nil
}

// ChangePassword 修改用户密码
func (as AuthService) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) error {
	// 检查确认密码是否匹配新密码（额外保险，即使已通过验证器）
	if req.NewPassword != req.ConfirmPassword {
		return utils.NewBusinessError(utils.PasswordMismatchCode, "")
	}

	// 新旧密码一样
	if req.NewPassword == req.OldPassword {
		return utils.NewBusinessError(utils.NewPasswordSameAsOldCode, "")
	}

	// 确定 Identifier 类型（用户名、邮箱或手机号）
	identifierType := utils.IdentifyType(req.Identifier)
	var admin *entity.Users
	var err error

	switch identifierType {
	case "email":
		admin, err = as.dao.GetByFields(ctx, "", req.Identifier, "")
	case "phone":
		admin, err = as.dao.GetByFields(ctx, "", "", req.Identifier)
	default:
		admin, err = as.dao.GetByFields(ctx, req.Identifier, "", "")
	}

	// 检查用户是否存在
	if err != nil {
		logger.Logger.Errorf("[ChangePassword] Error fetching user by %s: %v", identifierType, err)
		return utils.NewBusinessError(utils.InternalError, "")
	}
	if admin == nil {
		return utils.NewBusinessError(utils.UserNotFoundCode, "")
	}

	// 验证旧密码是否正确
	isMatch, err := utils.VerifyPassword(req.OldPassword, admin.Password)
	if err != nil {
		// 如果发生了错误（非匹配错误），记录日志并返回
		logger.Logger.Errorf("[ChangePassword] Error verifying password: %v", err)
		return utils.NewBusinessError(utils.InternalError, "")
	}
	if !isMatch {
		// 如果密码不匹配，返回无效凭证错误
		return utils.NewBusinessError(utils.OldPasswordIncorrectCode, "")
	}

	// 加密新密码
	hashedPassword, err := utils.EncryptPassword(req.NewPassword)
	if err != nil {
		logger.Logger.Errorf("[ChangePassword] Error encrypting new password: %v", err)
		return utils.NewBusinessError(utils.PasswordGenerationFailedCode, "")
	}

	// 更新密码
	updates := map[string]interface{}{
		entity.UsersColumns.Password:  hashedPassword,
		entity.UsersColumns.UpdatedAt: time.Now(),
	}
	if err := as.dao.Update(ctx, admin.ID, updates); err != nil {
		logger.Logger.Errorf("[ChangePassword] Error updating password for user %s: %v", admin.ID, err)
		return utils.NewBusinessError(utils.PasswordChangeFailedCode, "")
	}

	return nil
}
