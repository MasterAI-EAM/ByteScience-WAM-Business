package v1

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/internal/api/data"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitDataRouter(routerGroup *gin.RouterGroup) {
	secret := conf.GlobalConf.Jwt.AccessSecret

	dataGroup := routerGroup.Group("/data", middleware.JWTAuth(secret))
	{
		sampleApi := data.NewSampleApi()
		utils.RegisterRoute(dataGroup, http.MethodGet, "/sample", sampleApi.List)
	}

}
