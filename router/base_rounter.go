package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitBaseRoutes(r *gin.RouterGroup) gin.IRoutes {
	base := r.Group("/base")
	fmt.Println(base)
	base.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return base
}
