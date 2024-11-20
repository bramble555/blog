package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitChatGroupRouters(r *gin.RouterGroup) gin.IRoutes {
	r.GET("/chat_groups", controller.GetChatGroupHandler)
	return r
}
