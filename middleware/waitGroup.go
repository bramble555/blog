package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
)

func WaitGroupMiddleware(wg *sync.WaitGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		wg.Add(1)       // 增加 WaitGroup 计数
		defer wg.Done() // 请求处理完成后调用 Done
		c.Next()        // 继续处理请求
	}
}
