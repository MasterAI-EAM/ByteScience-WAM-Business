package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CorsMiddleware 创建CORS跨域中间件
func CorsMiddleware(allowOrigins string, allowMethods string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置跨域相关头部
		c.Header("Access-Control-Allow-Origin", allowOrigins)
		c.Header("Access-Control-Allow-Methods", allowMethods)
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
