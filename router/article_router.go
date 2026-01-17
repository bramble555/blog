package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitArticleRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/articles", middleware.JWTAuthorMiddleware(), controller.UploadArticleHandler)
	r.GET("/articles", controller.GetArticlesListHandler)
	r.GET("/articles/:sn", middleware.JWTOptionalMiddleware(), controller.GetArticlesDetailHandler)
	r.GET("/articles/calendar", middleware.JWTAuthorMiddleware(), controller.GetArticlesCalendarHandler)
	r.PUT("articles/:sn", middleware.JWTAdminMiddleware(), controller.UpdateArticlesHandler)
	r.DELETE("articles", middleware.JWTAdminMiddleware(), controller.DeleteArticlesListHandler)
	r.POST("articles/digg", middleware.JWTAuthorMiddleware(), controller.PostArticleDigHandler)
	r.POST("articles/collects", middleware.JWTAuthorMiddleware(), controller.PostArticleCollectHandler)
	r.GET("articles/collects", middleware.JWTAuthorMiddleware(), controller.GetArticleCollectHandler)
	r.DELETE("articles/collects", middleware.JWTAuthorMiddleware(), controller.DeleteArticleCollectHandler)
	return r
}
