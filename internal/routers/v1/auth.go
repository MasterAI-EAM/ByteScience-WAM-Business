package v1

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/internal/api/auth"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitAuthRouter(routerGroup *gin.RouterGroup) {
	secret := conf.GlobalConf.Jwt.AccessSecret

	authApi := auth.NewAuthApi()
	{
		utils.RegisterRoute(routerGroup, http.MethodPost, "/login", authApi.Login)
		utils.RegisterRoute(routerGroup, http.MethodPut, "/changPassword", authApi.ChangPassword)
	}

	// 假设您有一个路由组：/auth
	authGroup := routerGroup.Group("/search", middleware.JWTAuth(secret))
	{
		utils.RegisterRoute(authGroup, http.MethodPut, "/bs", authApi.ChangPassword)
	}

}
