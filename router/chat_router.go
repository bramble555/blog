package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitChatGroupRouters(r *gin.RouterGroup) gin.IRoutes {
	r.GET("/chat_groups", middleware.JWTAuthorMiddleware(), controller.GetChatGroupHandler)
	r.GET("/chat_groups_records", middleware.JWTAuthorMiddleware(), controller.GetChatGroupRecordsHandler)
	r.POST("/chat_groups_images", middleware.JWTAuthorMiddleware(), controller.UploadChatGroupImageHandler)
	return r
}
