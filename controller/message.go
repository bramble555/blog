package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

func SendMessageHandler(c *gin.Context) {
	var p model.ParamSendMessage
	if err := c.ShouldBindJSON(&p); err != nil {
		global.Log.Errorf("controller SendMessageHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	data, err := logic.SendMessage(claims.SN, &p)
	if err != nil {
		global.Log.Errorf("controller SendMessageHandler logic.SendMessage err:%s\n", err.Error())
		ResponseErrorWithErr(c, err, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

// GetMyMessagesListHandler 获取我的消息列表, 根据用户ID获取用户的消息列表
func GetMyMessagesListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetMyMessagesListHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	data, err := logic.GetMyMessagesList(claims.SN, pl)
	if err != nil {
		global.Log.Errorf("controller GetMyMessagesHandler logic.GetMyMessages err:%s\n", err.Error())
		ResponseErrorWithErr(c, err, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

// GetSentMessagesListHandler 获取我已发送的消息列表
func GetSentMessagesListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetSentMessagesListHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	data, err := logic.GetSentMessagesList(claims.SN, pl)
	if err != nil {
		global.Log.Errorf("controller GetSentMessagesListHandler logic.GetSentMessagesList err:%s\n", err.Error())
		ResponseErrorWithErr(c, err, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

// GetMessagesAllListHandler 获取所有消息
func GetMessagesAllListHandler(c *gin.Context) {
	var p model.ParamList
	if err := c.ShouldBindQuery(&p); err != nil {
		// Ignore error
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}

	data, err := logic.GetMessagesAllList(&p)
	if err != nil {
		global.Log.Errorf("controller GetMessagesAllHandler logic.GetMessagesAll err:%s\n", err.Error())
		ResponseErrorWithErr(c, err, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

// BroadcastMessageHandler 广播消息, 管理员发送消息给所有用户
func BroadcastMessageHandler(c *gin.Context) {
	var p model.ParamBroadcastMessage
	if err := c.ShouldBindJSON(&p); err != nil {
		global.Log.Errorf("controller BroadcastMessageHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	data, err := logic.BroadcastMessage(claims.SN, &p)
	if err != nil {
		global.Log.Errorf("controller BroadcastMessageHandler logic.BroadcastMessage err:%s\n", err.Error())
		ResponseErrorWithErr(c, err, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

// ReadMessageHandler 读取消息, 标记消息为已读
func ReadMessageHandler(c *gin.Context) {
	var p model.ParamSN
	if err := c.ShouldBindJSON(&p); err != nil {
		global.Log.Errorf("controller ReadMessageHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	data, err := logic.ReadMessage(claims.SN, p.SN)
	if err != nil {
		global.Log.Errorf("controller ReadMessageHandler logic.ReadMessage err:%s\n", err.Error())
		ResponseErrorWithErr(c, err, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
