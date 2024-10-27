package router

import (
	"github.com/bramble555/blog/controller"
	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("/menus", controller.UploadMenuHandler)
	// 分页查询
	r.GET("/menus", controller.GetMenuListHandler)
	r.PUT("/menus/:id", controller.UpdateMenuHandler)
	r.DELETE("/menus", controller.DeleteMenuListHander)
	return r
}
