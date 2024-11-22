package auth

type LoginRequest struct {
	// Identifier 用户标识（用户名|手机号|邮箱），必填，长度限制
	Identifier string `json:"identifier" validate:"required,min=3,max=128" example:"user1@example.com"`
	// Password 密码，必填，长度限制
	Password string `json:"password"  validate:"required,min=6,max=32" example:"password123"`
}

type LoginResponse struct {
	// Token 登陆凭证
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzIwOTkwMTQsImlhdCI6MTczMjA5OTAxNCwidXNlcklkIjoiMmE1ZWVkNDItMjVhMy00MGJlLTlmY2QtNjEzMmJlYzgzNTE3In0.YyrvQS66uYNVtCKKi7rm7xqJrCIFSq12SXCJcqAxKso"`
}

type ChangePasswordRequest struct {
	// Identifier 用户标识（用户名|手机号|邮箱），必填，长度限制
	Identifier string `json:"identifier" validate:"required,min=3,max=128" example:"user1@example.com"`
	// OldPassword 旧密码，必填，长度限制
	OldPassword string `json:"oldPassword" validate:"required,min=6,max=32" example:"oldpassword123"`
	// NewPassword 新密码，必填，长度限制
	NewPassword string `json:"newPassword" validate:"required,min=6,max=32" example:"newpassword123"`
	// ConfirmPassword 确认新密码，必填，必须与新密码一致
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword" example:"newpassword123"`
}
