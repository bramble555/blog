package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitBannerRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/images", middleware.JWTAuthorMiddleware(), controller.UploadBannersHandler)
	r.GET("/images", middleware.JWTAuthorMiddleware(), controller.GetBannerListHandler)
	r.DELETE("/images", middleware.JWTAdminMiddleware(), controller.DeleteBannerListHandler)
	return r
}
