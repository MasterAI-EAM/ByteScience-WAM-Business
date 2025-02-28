package v1

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/internal/api/ai"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitAiRouter(routerGroup *gin.RouterGroup) {
	secret := conf.GlobalConf.Jwt.AccessSecret

	aiGroup := routerGroup.Group("/ai", middleware.JWTAuth(secret))
	{
		predictApi := ai.NewPredictApi()
		utils.RegisterRoute(aiGroup, http.MethodPost, "/predict/forwardDirection", predictApi.ForwardDirection)
	}

}
