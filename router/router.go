package router

import (
	"github.com/gin-gonic/gin"
)

func InitRounter(mode string) *gin.Engine {
	// 如果是发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}
	r := gin.Default()
	r.GET("ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	apiGroup := r.Group("/api")
	InitBaseRoutes(apiGroup)
	InitBannerRoutes(apiGroup)
	InitAdvertRoutes(apiGroup)
	InitMenuRoutes(apiGroup)
	InitUserRoutes(apiGroup)
	InitMessageRoutes(apiGroup)
	InitArticleRoutes(apiGroup)
	InitTagRoutes(apiGroup)
	InitCommentRoutes(apiGroup)
	return r
}
