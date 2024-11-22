package utils

// 错误码定义
const (
	Success       = 0   // 成功
	BadRequest    = 400 // 请求错误
	InternalError = 500 // 服务器内部错误

	// 用户模块
	UserAlreadyExistsCode      = 1001 // 用户已存在
	UserNotFoundCode           = 1002 // 用户未找到
	UserInvalidCredentialsCode = 1003 // 用户凭证无效
	UsernameAlreadyExistsCode  = 1004 // 用户名已存在
	EmailAlreadyExistsCode     = 1005 // 邮箱已存在
	PhoneAlreadyExistsCode     = 1006 // 手机号已存在

	// 管理员模块
	AdminAlreadyExistsCode         = 1101 // 管理员已存在
	AdminNotFoundCode              = 1102 // 管理员未找到
	AdminInvalidIDCode             = 1103 // 管理员ID无效
	AdminUnauthorizedCode          = 1106 // 管理员权限不足
	AdminUsernameAlreadyExistsCode = 1107 // 管理员用户名已存在
	AdminPhoneAlreadyExistsCode    = 1108 // 管理员手机已存在
	AdminEmailAlreadyExistsCode    = 1109 // 管理员邮箱已存在

	// 权限系统模块
	PermissionDeniedCode            = 1201 // 权限被拒绝
	InsufficientRolePermissionsCode = 1202 // 角色权限不足
	InvalidTokenCode                = 1203 // 无效的令牌
	TokenExpiredCode                = 1204 // 令牌已过期
	AccessForbiddenCode             = 1205 // 禁止访问
	RoleNameAlreadyExistsCode       = 1206 // 角色名称已存在
	RoleNotFoundCode                = 1207 // 角色未找到
	RoleAssignmentFailedCode        = 1208 // 角色分配失败

	// 密码相关
	PasswordIncorrectCode        = 1301 // 密码不正确
	PasswordTooWeakCode          = 1302 // 密码过于简单
	PasswordResetFailedCode      = 1303 // 密码重置失败
	PasswordChangeFailedCode     = 1304 // 密码更改失败
	PasswordMismatchCode         = 1305 // 新密码与确认密码不匹配
	OldPasswordIncorrectCode     = 1306 // 旧密码不正确
	PasswordGenerationFailedCode = 1307 // 密码生成失败
	NewPasswordSameAsOldCode     = 1308 // 新密码与旧密码相同

	// 接口错误
	AdminInsertFailedCode       = 2001 // 插入管理员失败
	AdminUpdateFailedCode       = 2002 // 更新管理员信息失败
	AdminDeleteFailedCode       = 2003 // 删除管理员失败
	AdminQueryListFailedCode    = 2004 // 查询管理员列表失败
	RoleInsertFailedCode        = 2005 // 插入角色失败
	RoleUpdateFailedCode        = 2006 // 更新角色失败
	RoleDeleteFailedCode        = 2007 // 删除角色失败
	RoleQueryListFailedCode     = 2008 // 查询角色列表失败
	UserDeleteFailedCode        = 2009 // 用户删除失败
	UserQueryFailedCode         = 2010 // 用户查询失败
	UserQueryListFailedCode     = 2011 // 用户列表查询失败
	UserInsertFailedCode        = 2012 // 用户插入失败
	UserConflictCheckFailedCode = 2013 // 用户冲突检测失败
	UserUpdateFailedCode        = 2014 // 用户更新失败
)

// ErrorMessages 错误信息映射
var ErrorMessages = map[int]string{
	Success:       "success",
	BadRequest:    "Invalid Request Parameters",
	InternalError: "Internal Server Error",

	// 用户模块
	UserAlreadyExistsCode:      "User already exists",
	UserNotFoundCode:           "User not found",
	UserInvalidCredentialsCode: "Invalid credentials",
	UsernameAlreadyExistsCode:  "Username already exists",
	EmailAlreadyExistsCode:     "Email already exists",
	PhoneAlreadyExistsCode:     "Phone number already exists",

	// 管理员模块
	AdminAlreadyExistsCode:         "Admin already exists",
	AdminNotFoundCode:              "Admin not found",
	AdminInvalidIDCode:             "Invalid Admin ID",
	AdminUnauthorizedCode:          "Admin lacks sufficient permissions",
	AdminUsernameAlreadyExistsCode: "Admin username already exists",
	AdminEmailAlreadyExistsCode:    "Admin email already exists",
	AdminPhoneAlreadyExistsCode:    "Admin phone already exists",
	NewPasswordSameAsOldCode:       "Admin newPassword same as old",

	// 权限系统模块
	PermissionDeniedCode:            "Permission denied",
	InsufficientRolePermissionsCode: "Insufficient role permissions",
	InvalidTokenCode:                "Invalid token",
	TokenExpiredCode:                "Token expired",
	AccessForbiddenCode:             "Access forbidden",
	RoleNotFoundCode:                "Role not found",
	RoleAssignmentFailedCode:        "Failed to assign role",
	RoleNameAlreadyExistsCode:       "Role name already exists",

	// 密码相关
	PasswordIncorrectCode:        "Incorrect password",
	PasswordTooWeakCode:          "Password is too weak",
	PasswordResetFailedCode:      "Password reset failed",
	PasswordChangeFailedCode:     "Password change failed",
	PasswordMismatchCode:         "Passwords do not match",
	OldPasswordIncorrectCode:     "Old password is incorrect",
	PasswordGenerationFailedCode: "Failed to generate password",

	// 接口错误
	AdminInsertFailedCode:       "Failed to insert admin",
	AdminUpdateFailedCode:       "Failed to update admin",
	AdminDeleteFailedCode:       "Failed to delete admin",
	AdminQueryListFailedCode:    "Failed to query admin list",
	RoleInsertFailedCode:        "Failed to insert role",
	RoleUpdateFailedCode:        "Failed to update role",
	RoleDeleteFailedCode:        "Failed to delete role",
	RoleQueryListFailedCode:     "Failed to query role list",
	UserConflictCheckFailedCode: "User conflict check failed",
	UserUpdateFailedCode:        "Failed to update user",
	UserDeleteFailedCode:        "Failed to delete user",
	UserQueryFailedCode:         "Failed to query user",
	UserQueryListFailedCode:     "Failed to query user list",
	UserInsertFailedCode:        "Failed to insert user",
}