package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitTagRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/tags", controller.CreateTagsHandle)
	r.GET("/tags", controller.GetTagsListHandler)
	r.DELETE("/tags", controller.DeleteTagsListHandler)
	return r
}
