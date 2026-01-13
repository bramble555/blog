package message

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func SendMessage(pm *model.ParamMessage) (string, error) {
	// 创建消息
	message := model.MessageModel{
		SendUserSN: pm.SendUserSN,
		RevUserSN:  pm.RevUserSN,
		Content:    pm.Content,
		IsRead:     false,
	}
	// 插入数据
	err := global.DB.Create(&message).Error
	if err != nil {
		global.Log.Errorf("Send Message Create Error: %s", err)
		return "", err
	}
	return "消息发送成功", nil
}
