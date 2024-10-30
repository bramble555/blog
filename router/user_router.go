package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/user", controller.EmailLoginHandler)
	return r
}
