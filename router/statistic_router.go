package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitStatisticRoutes(r *gin.RouterGroup) {
	statisticRouter := r.Group("data")
	statisticRouter.Use(middleware.JWTAuthorMiddleware())
	{
		statisticRouter.GET("login_data", controller.GetUserLoginHandler)
		statisticRouter.GET("sum", controller.GetDataSumHandler)
	}
}
