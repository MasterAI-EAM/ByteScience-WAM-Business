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
		utils.RegisterRoute(dataGroup, http.MethodPut, "/sample", sampleApi.Edit)
		utils.RegisterRoute(dataGroup, http.MethodDelete, "/sample", sampleApi.Delete)
		utils.RegisterRoute(dataGroup, http.MethodPost, "/sample/import", sampleApi.Import)

		inferenceApi := data.NewInferenceApi()
		utils.RegisterRoute(dataGroup, http.MethodPost, "/inference/prediction", inferenceApi.Prediction)
	}

}
