package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

func SendMessageHandler(c *gin.Context) {
	pm := model.ParamMessage{}
	err := c.ShouldBindJSON(&pm)
	if err != nil {
		global.Log.Errorf("controller ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.SendMessage(&pm)
	if err != nil {
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
