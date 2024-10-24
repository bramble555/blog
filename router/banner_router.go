package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitBannerRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/images", controller.UploadBannerHandler)
	// 分页查询
	r.GET("/images", controller.GetBannerListHandler)
	// 不分页，并且不包含 create_time 等信息
	r.GET("/images_detail", controller.GetBannerDetailHandler)
	r.DELETE("/images", controller.DeleteBannerListHander)
	return r
}
