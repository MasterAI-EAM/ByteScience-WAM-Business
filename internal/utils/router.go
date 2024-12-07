package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterRoute 封装了请求参数绑定和统一响应，支持具体返回类型，并支持中间件
func RegisterRoute[T any, R any](group *gin.RouterGroup, method, path string, handlerFunc func(ctx *gin.Context, req *T) (R, error), middlewares ...gin.HandlerFunc) {
	group.Handle(method, path, append(middlewares, func(ctx *gin.Context) {
		// 执行请求参数绑定和校验
		req, err := bindAndValidate[T](ctx, method)
		if err != nil {
			SendResponse(ctx, http.StatusBadRequest, ErrorResponse(BadRequest, err.Error()))
			return
		}

		// 调用实际的处理函数
		res, err := handlerFunc(ctx, req)
		if err != nil {
			// 检查错误类型，如果是 BusinessError，则直接响应
			if businessErr, ok := err.(*BusinessError); ok {
				SendResponse(ctx, http.StatusBadRequest, ErrorResponse(businessErr.Code, businessErr.Message))
				return
			}

			// 非业务错误返回 500 状态码
			SendResponse(ctx, http.StatusInternalServerError, ErrorResponse(InternalError, err.Error()))
			return
		}

		// 成功响应
		SendResponse(ctx, http.StatusOK, SuccessResponse(res))
	})...)
}
