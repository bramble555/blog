package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func CORS() gin.HandlerFunc {
	// 创建 CORS 处理
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8080", "null"},
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "token"},
		AllowCredentials: true,
		MaxAge:           int(1 * time.Hour / time.Second), // 转换为秒
	})

	return func(c *gin.Context) {
		// 这里创建一个 http.HandlerFunc，传递给 corsHandler.ServeHTTP
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next() // 调用 Gin 的下一个中间件或请求处理函数
		})

		// 使用 CORS 处理器来处理请求
		corsHandler.ServeHTTP(c.Writer, c.Request, next)
	}
}
