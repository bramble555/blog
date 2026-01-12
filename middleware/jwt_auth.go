package middleware

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model/ctype"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

// JWTAuthorMiddleware 登录权限认证
func JWTAuthorMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		ok := redis.CheckLogout(authHeader)
		if ok {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 解析 token
		claims, err := pkg.ParseToken(authHeader)
		if err != nil {
			global.Log.Errorf("ParseToken failed: %v, token: [%s]", err, authHeader)
			controller.ResponseError(c, controller.CodeInvalidAuth)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("claims", claims)
		c.Next()
	}
}

// JWTAdminMiddleware 管理员权限认证
func JWTAdminMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		ok := redis.CheckLogout(authHeader)
		if ok {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 解析 token
		claims, err := pkg.ParseToken(authHeader)
		if err != nil {
			global.Log.Errorf("Admin ParseToken failed: %v, token: [%s]", err, authHeader)
			controller.ResponseError(c, controller.CodeInvalidAuth)
			c.Abort()
			return
		}
		if claims.Role != uint(ctype.PermissionAdmin) {
			controller.ResponseErrorWithData(c, controller.CodeInvalidAuth, "您不是管理员身份")
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("claims", claims)
		c.Next()
	}
}
