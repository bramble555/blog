package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitArtilceRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/articles", middleware.JWTAuthorMiddleware(), controller.UploadArticlesHandler)
	r.GET("/articles", controller.GetArticlesListHandler)
	return r
}
