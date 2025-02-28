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

		recipeApi := data.NewRecipeApi()
		utils.RegisterRoute(dataGroup, http.MethodPost, "/recipe", recipeApi.Add)
		utils.RegisterRoute(dataGroup, http.MethodGet, "/recipe", recipeApi.List)
		utils.RegisterRoute(dataGroup, http.MethodDelete, "/recipe", recipeApi.Delete)
		utils.RegisterRoute(dataGroup, http.MethodPut, "/recipe", recipeApi.Edit)
		utils.RegisterRoute(dataGroup, http.MethodGet, "/recipe/info", recipeApi.Info)
		utils.RegisterRoute(dataGroup, http.MethodGet, "/recipe/form/list", recipeApi.FormList)
	}

}
