package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitArticleRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/articles", middleware.JWTAuthorMiddleware(), controller.UploadArticlesHandler)
	r.GET("/articles", controller.GetArticlesListHandler)
	r.GET("/articles/:id", controller.GetArticlesDetailHandler)
	r.GET("/articles/calendar", controller.GetArticlesCalendarHandler)
	r.GET("/articles/tags", controller.GetArticlesTagsListHandler)
	r.PUT("articles/:id", controller.UpdateArticlesHandler)
	return r
}
