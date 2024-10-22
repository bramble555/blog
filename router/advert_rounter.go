package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitAdvertRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/", controller.CreateAdvertHandle)
	return r
}
