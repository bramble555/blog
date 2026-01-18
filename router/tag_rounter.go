package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitTagRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/tags", middleware.JWTAuthorMiddleware(), controller.CreateTagsHandle)
	r.GET("/tags", middleware.JWTAuthorMiddleware(), controller.GetTagsHandler)
	r.DELETE("/tags", middleware.JWTAdminMiddleware(), controller.DeleteTagsListHandler)
	return r
}
