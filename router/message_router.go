package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitMessageRoutes(r *gin.RouterGroup) gin.IRoutes {

	r.POST("messages", middleware.JWTAuthorMiddleware(), controller.SendMessageHandler)
	// 	r.GET("messages_all", middleware.JwtAdmin(), controller.MessageListAllView)
	// 	r.GET("messages", middleware.JwtAuth(), controller.MessageListView)
	// 	r.GET("messages_record", middleware.JwtAuth(), controller.MessageRecordView)
	return r
}
