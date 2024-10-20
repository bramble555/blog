package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRounter(mode string) *gin.Engine {
	// 如果是发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}
	r := gin.Default()
	r.GET("ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
	
	return r
}
