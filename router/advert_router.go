package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitAdvertRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/adverts", middleware.JWTAuthorMiddleware(), controller.UploadAdvertImagesHandler)
	r.GET("/adverts", controller.GetAdvertListHandler)
	r.DELETE("/adverts", middleware.JWTAdminMiddleware(), controller.DeleteAdvertListHandler)
	r.PUT("/adverts", middleware.JWTAdminMiddleware(), controller.UpdateAdvertShowHandler)
	return r
}
