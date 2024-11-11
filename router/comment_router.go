package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitCommentRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/comments", middleware.JWTAuthorMiddleware(), controller.PostArticleCommentsHandler)
	r.GET("/comments", controller.GetArticleCommentsHandler)
	r.POST("/comments/digg", middleware.JWTAuthorMiddleware(), controller.PostArticleCommentsDiggHandler)
	r.DELETE("/comments", middleware.JWTAuthorMiddleware(), controller.DeleteArticleCommentsHandler)
	return r
}
