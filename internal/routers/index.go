package routers

import (
	"ByteScience-WAM-Business/docs" // 导入 Swagger 生成的文档
	"ByteScience-WAM-Business/internal/routers/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register(router *gin.Engine) {
	// 注册swagger路由
	docs.SwaggerInfo.BasePath = "/v1"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 加载v1的路由
	v1.LoadRouters(router)
}
