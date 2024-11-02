package message

import (
	"github.com/bramble555/blog/dao/mysql"
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

func MessageListAll(pl *model.ParamList) ([]model.RespondMessage, error) {
	return mysql.GetTableList[model.RespondMessage]("message_models", pl, "")
}
func MessageList(id uint) ([]model.RespondMessage, error) {
	// 分组，相当于 temp
	var messageGroup = model.RespondMessageGroup{}
	// 数据库里面的数据在 messageList 里面
	var messageList []model.MessageModel
	// 要返回的 res
	// 内容是最新内容，并且包含消息创建时间
	var messageRespond []model.RespondMessage

	// 查找和本人有关的
	global.DB.Order("create_time asc").
		Find(&messageList, "send_user_id = ? or rev_user_id  = ?", id, id)
	for _, m := range messageList {
		// 判断是一个组的条件
		// 1 2 和 2 1；1 3 和 3 1 是一组
		message := model.RespondMessage{
			SendUserID:     m.SendUserID,
			SendUsername:   m.SendUsername,
			SendUserAvatar: m.SendUserAvatar,
			RevUserID:      m.RevUserID,
			RevUsername:    m.RevUsername,
			RevUserAvatar:  m.RevUserAvatar,
			Content:        m.Content,
			MessageCount:   0,
		}

		uniqueIDSum := m.SendUserID + m.RevUserID
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
func MessageRecord(myID, recordID uint) ([]model.MessageModel, error) {
	// 查询之后的 _messageList 列表，然后需要分组进行返回
	var _messageList []model.MessageModel
	var messageList = make([]model.MessageModel, 0)
	// 先查找与自己相关联的聊天记录，然后分组,然后再寻找 recordID
	err := global.DB.Order("create_time asc").Find(&_messageList, "send_user_id = ? or rev_user_id = ?", myID, myID).Error
	if err != nil {
		global.Log.Errorf("message MessageRecord err:%s\n", err.Error())
		return nil, err
	}
	for _, model := range _messageList {
		if model.RevUserID == recordID || model.SendUserID == recordID {
			messageList = append(messageList, model)
		}
	}
	return messageList, nil
}
