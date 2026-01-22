package router

import (
	"sync"

	"github.com/bramble555/blog/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

var store = cookie.NewStore([]byte("bbbbbb"))

func InitRouter(mode string, wg *sync.WaitGroup) *gin.Engine {
	// 如果是发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}
	r := gin.Default()
	r.GET("ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.Use(middleware.CORS())
	// 如果强制退出，必须要把正在执行的任务处理完才退出
	r.Use(middleware.WaitGroupMiddleware(wg))
	r.Use(middleware.RateLimitMiddleware())
	if store == nil {
		store = cookie.NewStore([]byte("bbbbbb"))
	}
	r.Use(sessions.Sessions("sessiodddnd", store))
	r.Static("/uploads", "./uploads")
	apiGroup := r.Group("/api")
	InitUserRoutes(apiGroup)
	InitBannerRoutes(apiGroup)
	InitAdvertRoutes(apiGroup)
	InitTagRoutes(apiGroup)
	InitArticleRoutes(apiGroup)
	InitMessageRoutes(apiGroup)
	InitCommentRoutes(apiGroup)
	InitChatGroupRouters(apiGroup)
	InitMetricsRoutes(apiGroup)
	InitStatisticRoutes(apiGroup)
	return r
}
