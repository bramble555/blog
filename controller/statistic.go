package controller

import (
	"github.com/bramble555/blog/logic"
	"github.com/gin-gonic/gin"
)

func GetUserLoginHandler(c *gin.Context) {
	data, err := logic.GetUserLoginData()
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetDataSumHandler(c *gin.Context) {
	data, err := logic.GetDataSum()
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
