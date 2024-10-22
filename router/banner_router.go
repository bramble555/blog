package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitBannerRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/images", controller.UploadBannerHandler)
	r.GET("/images", controller.GetBannerListHandler)
	r.DELETE("/images", controller.DeleteBannerListHander)
	return r
}
