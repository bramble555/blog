package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitMessageRoutes(r *gin.RouterGroup) gin.IRoutes {

	r.POST("messages", middleware.JWTAuthorMiddleware(), controller.SendMessageHandler)
	r.GET("messages_all", middleware.JWTAdminMiddleware(), controller.MessageListAllHandler)
	r.GET("messages", middleware.JWTAuthorMiddleware(), controller.MessageListHandler)
	r.GET("messages_record", middleware.JWTAuthorMiddleware(), controller.MessageRecordHandler)
	return r
}
