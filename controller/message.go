package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

func SendMessageHandler(c *gin.Context) {
	pm := model.ParamMessage{}
	err := c.ShouldBindJSON(&pm)
	if err != nil {
		global.Log.Errorf("controller SendMessageHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.SendMessage(&pm)
	if err != nil {
		global.Log.Errorf("controller SendMessageHandler logic.SendMessage err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

// MessageListAllHandler admin 以列表形式查看全部信息
func MessageListAllHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller MessageListAllHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.MessageListAll(pl)
	if err != nil {
		global.Log.Errorf("controller MessageListAllHandler logic.MessageListAll err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

// MessageListAllHandler user 以组形式查看与他人有关的信息
// 1 3  3 1 是一个组
func MessageListHandler(c *gin.Context) {
	_cliams, _ := c.Get("claims")
	cliams := _cliams.(*pkg.MyClaims)
	data, err := logic.MessageList(cliams.SN)
	if err != nil {
		global.Log.Errorf("controller MessageListHandler logic.MessageList err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
func MessageRecordHandler(c *gin.Context) {
	pr := model.ParamRecordSN{}
	err := c.ShouldBindJSON(&pr)
	if err != nil {
		global.Log.Errorf("controller MessageRecordHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	data, err := logic.MessageRecord(claims.SN, pr.UserSN)
	if err != nil {
		global.Log.Errorf("controller MessageRecordHandler logic.MessageRecord err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	// 点开消息，里面的每一条消息都从未读变为已读
	ResponseSucceed(c, data)
}
