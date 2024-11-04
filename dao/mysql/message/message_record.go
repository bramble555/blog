package message

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

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
