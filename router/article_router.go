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
	r.DELETE("articles", controller.DeleteArticlesListHandler)
	r.POST("articles/digg", controller.PostArticleDigHandler)
	r.POST("articles/collects", middleware.JWTAuthorMiddleware(), controller.PostArticleCollectHandler)
	r.GET("articles/collects", middleware.JWTAuthorMiddleware(), controller.GetArticleCollectHandler)
	r.DELETE("articles/collects", middleware.JWTAuthorMiddleware(), controller.DeleteArticleCollectHandler)
	return r
}
