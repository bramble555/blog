package router

import (
	"sync"

	"github.com/bramble555/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(mode string, wg *sync.WaitGroup) *gin.Engine {
	// 如果是发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}
	r := gin.Default()
	r.GET("ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	// 如果强制退出，必须要把正在执行的任务处理完才退出
	r.Use(middleware.WaitGroupMiddleware(wg))
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
