package utils

// 错误码定义
const (
	Success              = 0   // 成功
	BadRequest           = 400 // 请求错误
	InternalError        = 500 // 服务器内部错误
	ExternalRequestError = 600 // 外部请求错误

	// 用户模块
	DatabaseErrorCode          = 1000 // 数据库错误
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

	// 文件相关
	FileTooLargeCode        = 1401 // 文件大小超出限制
	DuplicateFileImportCode = 1402 // 文件重复导入

	// 接口处理
	FileParsingFailedCode = 3001 // 文件解析失败

	// 业务错误

	MaterialProportionSumNot100Code   = 4001 // 材料组占比总和不是100%
	MaterialGroupProportionNot100Code = 4002 // 实验步骤中材料组的占比不是100%
	ExperimentDoesNotExistCode        = 4003 // 实验不存在
	DuplicateRecipeFormatCode         = 4004 // 配方重复
	DuplicateExperimentFormatCode     = 4005 // 实验重复
	RecipeDoesNotExistCode            = 4006 // 配方不存在
	RecipeDeletionIsNotAllowedCode    = 4007 // 不允许删除配方

)

// ErrorMessages 错误信息映射
var ErrorMessages = map[int]string{
	Success:       "success",
	BadRequest:    "Invalid Request Parameters",
	InternalError: "Internal Server Error",

	// 用户模块
	DatabaseErrorCode:          "Database Error",
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

	// 文件相关
	FileTooLargeCode:        "File too large",
	DuplicateFileImportCode: "Duplicate file import",

	// 接口处理
	FileParsingFailedCode: "File parsing failed",

	// 业务错误
	MaterialProportionSumNot100Code:   "The sum of the material proportions in the material group is not 100%",
	MaterialGroupProportionNot100Code: "The proportion of the material group is not 100%",
	ExperimentDoesNotExistCode:        "Experiment does not exist",
	RecipeDoesNotExistCode:            "Recipe does not exist",
	DuplicateRecipeFormatCode:         "This recipe is repeated with '%s'",
	DuplicateExperimentFormatCode:     "This experiment is repeated with '%s'",
	RecipeDeletionIsNotAllowedCode:    "The recipe is in use and cannot be deleted.",
}
