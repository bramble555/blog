package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) gin.IRoutes {

	// 配置路由
	r.POST("/username", controller.UsernameLoginHandler)
	r.GET("/users", middleware.JWTAuthorMiddleware(), controller.GetUserListHandler)
	r.PUT("/user_role", middleware.JWTAdminMiddleware(), controller.UpdateUserRoleHandler)
	r.PUT("/user_password", middleware.JWTAuthorMiddleware(), controller.UpdateUserPwdHandler)
	r.POST("/logout", middleware.JWTAuthorMiddleware(), controller.LogoutHandler)
	r.DELETE("/users", middleware.JWTAdminMiddleware(), controller.DeleteUserListHandler)
	r.POST("/user_bind_email", middleware.JWTAuthorMiddleware(), controller.PostBindEmailHandler)

	return r
}
