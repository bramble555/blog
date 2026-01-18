package message

import (
	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func MessageListAll(pl *model.ParamList) ([]model.RespondMessage, int64, error) {
	return mysql.GetTableList[model.RespondMessage]("message_models", pl, "")
}
func MessageList(sn int64) ([]model.RespondMessage, error) {
	// 分组，相当于 temp
	var messageGroup = model.RespondMessageGroup{}
	// 数据库里面的数据在 messageList 里面
	var messageList []model.MessageModel
	// 要返回的 res
	// 内容是最新内容，并且包含消息创建时间
	var messageRespond []model.RespondMessage

	// 查找和本人有关的
	global.DB.Order("create_time asc").
		Find(&messageList, "send_user_sn = ? or rev_user_sn  = ?", sn, sn)
	for _, m := range messageList {
		// 判断是一个组的条件
		// 1 2 和 2 1；1 3 和 3 1 是一组
		message := model.RespondMessage{
			SendUserSN:     m.SendUserSN,
			SendUsername:   m.SendUsername,
			SendUserAvatar: m.SendUserAvatar,
			RevUserSN:      m.RevUserSN,
			RevUsername:    m.RevUsername,
			RevUserAvatar:  m.RevUserAvatar,
			Content:        m.Content,
			MessageCount:   0,
		}

		uniqueIDSum := m.SendUserSN + m.RevUserSN
		val, exists := messageGroup[uniqueIDSum] // 直接获取 val 和 exists
		if !exists || val == nil {
			// 如果不存在或 val 为 nil，初始化 MessageCount
			message.MessageCount = 1
		} else {
			// 如果存在且不为 nil，更新 MessageCount
			message.MessageCount = val.MessageCount + 1
		}

		// 更新消息为最新消息
		messageGroup[uniqueIDSum] = &message
	}

	// 把分组内容添加到 messages
	for _, message := range messageGroup {
		messageRespond = append(messageRespond, *message)
	}
	return messageRespond, nil // 返回结果和 nil 错误
}
