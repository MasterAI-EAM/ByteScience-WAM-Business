package middleware

import (
	"ByteScience-WAM-Business/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth 中间件
func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			// 如果没有提供 Authorization header，返回 401 错误
			utils.SendResponse(ctx, 401, utils.ErrorResponse(utils.InvalidTokenCode, "Missing token"))
			ctx.Abort()
			return
		}

		// 检查 token 是否包含 "Bearer " 前缀
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			// 如果 token 格式不正确，返回 401 错误
			utils.SendResponse(ctx, 401, utils.ErrorResponse(utils.InvalidTokenCode, "Malformed token"))
			ctx.Abort()
			return
		}

		// 提取 token
		tokenString := authHeader[7:]

		// 解析 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 检查 token 是否使用了 HS256 签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, utils.NewBusinessError(utils.InvalidTokenCode)
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			// token 解析失败，可能是过期、无效等原因
			if err.Error() == "Token is expired" {
				// 如果 token 过期，返回业务错误
				utils.SendResponse(ctx, 401, utils.ErrorResponse(utils.TokenExpiredCode, "Token expired"))
			} else {
				// 其他错误
				utils.SendResponse(ctx, 401, utils.ErrorResponse(utils.InvalidTokenCode, "Invalid token"))
			}
			ctx.Abort()
			return
		}

		// token 验证成功，设置用户信息
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 设置用户信息到上下文中
			ctx.Set("userId", claims["userId"])
		} else {
			// token 无效，返回错误
			utils.SendResponse(ctx, 401, utils.ErrorResponse(utils.InvalidTokenCode, "Invalid token"))
			ctx.Abort()
			return
		}

		// 继续请求处理
		ctx.Next()
	}
}
