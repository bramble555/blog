package message

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func SendMessage(pm *model.ParamMessage, udl []*model.UserDetail) (string, error) {
	// 创建消息
	message := model.MessageModel{
		SendUserID:     pm.SendUserID,
		SendUsername:   udl[0].Username,
		SendUserAvatar: udl[0].Avatar,
		RevUserID:      pm.RevUserID,
		RevUsername:    udl[1].Username,
		RevUserAvatar:  udl[1].Avatar,
		Content:        pm.Content,
	}

	// 保存消息到数据库
	err := global.DB.Create(&message).Error
	if err != nil {
		global.Log.Errorf("发送消息失败: %s\n", err.Error())
		return "", err
	}
	return "发送消息成功", nil
}
