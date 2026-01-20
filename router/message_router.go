package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitMessageRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/message/send", middleware.JWTAuthorMiddleware(), controller.SendMessageHandler)
	r.POST("/message/broadcast", middleware.JWTAdminMiddleware(), controller.BroadcastMessageHandler)
	r.PUT("/message/read", middleware.JWTAuthorMiddleware(), controller.ReadMessageHandler)
	r.GET("/messages", middleware.JWTAuthorMiddleware(), controller.GetMyMessagesListHandler)
	r.GET("/messages_sent", middleware.JWTAuthorMiddleware(), controller.GetSentMessagesListHandler)
	r.GET("/messages_all", middleware.JWTAdminMiddleware(), controller.GetMessagesAllListHandler)
	return r
}
