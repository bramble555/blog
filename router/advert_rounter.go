package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitAdvertRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/adverts", controller.CreateAdvertHandle)
	r.GET("/adverts", controller.GetAdvertListHandler)
	r.DELETE("/adverts", controller.DeleteAdvertListHander)
	return r
}
