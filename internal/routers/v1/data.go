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
		experimentApi := data.NewExperimentApi()
		utils.RegisterRoute(dataGroup, http.MethodGet, "/experiment", experimentApi.List)
		utils.RegisterRoute(dataGroup, http.MethodPost, "/experiment", experimentApi.Add)
		utils.RegisterRoute(dataGroup, http.MethodPut, "/experiment", experimentApi.Edit)
		utils.RegisterRoute(dataGroup, http.MethodDelete, "/experiment", experimentApi.Delete)
		utils.RegisterRoute(dataGroup, http.MethodPost, "/experiment/import", experimentApi.Import)

		inferenceApi := data.NewInferenceApi()
		utils.RegisterRoute(dataGroup, http.MethodPost, "/inference/prediction", inferenceApi.Prediction)
	}

}
