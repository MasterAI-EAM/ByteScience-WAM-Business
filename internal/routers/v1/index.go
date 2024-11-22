package v1

import (
	"github.com/gin-gonic/gin"
)

func LoadRouters(router *gin.Engine) {
	v1Group := router.Group("/v1")
	InitAuthRouter(v1Group)
}
