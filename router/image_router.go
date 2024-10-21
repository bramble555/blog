package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitImageRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/images", controller.ImageHandler)
	return r
}
