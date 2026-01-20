package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/gin-gonic/gin"
)

func GetChatGroupHandler(c *gin.Context) {
	if err := logic.HandleChatGroupWS(c); err != nil {
		global.Log.Errorf("controller GetChatGroupHandler HandleChatGroupWS err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
}

func GetChatGroupRecordsHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetAdvertListHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	res, err := logic.GetChatGroupRecords(pl)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, res)
}

func UploadChatGroupImageHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	url, err := logic.UploadChatGroupImage(c, file)
	if err != nil {
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}
	ResponseSucceed(c, url)
}
