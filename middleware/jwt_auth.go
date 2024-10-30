package middleware

import (
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

func JWTAuthorMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带 Token 有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设 Token 放在 Header 的 token 中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 解析 token
		claims, err := pkg.ParseToken(authHeader)
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidAuth)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("claims", claims)
		c.Next() // 后续的处理函数可以用过 c.Get("claims") 来获取当前请求的用户信息
	}
}
