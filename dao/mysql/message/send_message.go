package message

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func SendMessage(pm *model.ParamMessage) (string, error) {
	err := global.DB.Create(&model.MessageModel{
		SendUserID: pm.SendUserID,
		RevUserID:  pm.RevUserID,
		Content:    pm.Content,
	}).Error
	if err != nil {
		global.Log.Errorf("message SendMessage err:%s\n", err.Error())
		return "", err
	}
	return "发送消息成功", nil
}
