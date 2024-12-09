package v1

import (
	"ByteScience-WAM-Business/internal/api/auth"
	"ByteScience-WAM-Business/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitAuthRouter(routerGroup *gin.RouterGroup) {
	authApi := auth.NewAuthApi()
	{
		utils.RegisterRoute(routerGroup, http.MethodPost, "/login", authApi.Login)
		utils.RegisterRoute(routerGroup, http.MethodPut, "/changPassword", authApi.ChangPassword)
	}

}
