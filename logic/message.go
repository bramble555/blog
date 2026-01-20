package logic

import (
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/message"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/model"
)

func SendMessage(userSN int64, p *model.ParamSendMessage) (string, error) {
	if p == nil || p.RevUserSN == 0 || p.Content == "" {
		return "", code.ErrorInvalidParam
	}

	// 检查发送者是否存在
	sender, err := user.GetUserDetailBySN(userSN)
	if err != nil {
		return "", err
	}

	// 检查接收者是否存在
	receiver, err := user.GetUserDetailBySN(p.RevUserSN)
	if err != nil {
		return "", err
	}

	// 创建消息
	msg := model.MessageModel{
		SendUserSN:     userSN,
		SendUsername:   sender.Username,
		SendUserAvatar: sender.Avatar,
		RevUserSN:      receiver.SN,
		RevUsername:    receiver.Username,
		RevUserAvatar:  receiver.Avatar,
		IsRead:         false,
		Content:        p.Content,
	}

	return message.CreateMessage(&msg)
}

func BroadcastMessage(userSN int64, p *model.ParamBroadcastMessage) (string, error) {
	if p == nil || p.Content == "" {
		return "", code.ErrorInvalidParam
	}

	// 检查发送者是否存在
	sender, err := user.GetUserDetailBySN(userSN)
	if err != nil {
		return "", err
	}

	// 获取所有用户
	users, _, err := user.GetUserList(nil)
	if err != nil {
		return "", err
	}

	msgs := make([]model.MessageModel, 0, len(users))
	for _, u := range users {
		msgs = append(msgs, model.MessageModel{
			SendUserSN:     userSN,
			SendUsername:   sender.Username,
			SendUserAvatar: sender.Avatar,
			RevUserSN:      u.SN,
			RevUsername:    u.Username,
			RevUserAvatar:  u.Avatar,
			IsRead:         false,
			Content:        p.Content,
		})
	}
	return message.CreateMessagesBatch(msgs)
}

func ReadMessage(userSN int64, sn int64) (string, error) {
	if sn == 0 {
		return "", code.ErrorInvalidParam
	}
	return message.UpdateMessageRead(sn, userSN)
}

func GetMyMessagesList(userSN int64, pl *model.ParamList) (*model.PageResult[model.MessageModel], error) {
	list, count, err := message.GetUserMessageList(pl, userSN)
	if err != nil {
		return nil, err
	}
	return &model.PageResult[model.MessageModel]{
		List:  list,
		Count: count,
	}, nil
}

func GetSentMessagesList(userSN int64, pl *model.ParamList) (*model.PageResult[model.MessageModel], error) {
	list, count, err := message.GetSentMessageList(pl, userSN)
	if err != nil {
		return nil, err
	}
	return &model.PageResult[model.MessageModel]{
		List:  list,
		Count: count,
	}, nil
}

func GetMessagesAllList(pl *model.ParamList) (*model.PageResult[model.MessageModel], error) {
	list, count, err := message.GetAllMessagesList(pl)
	if err != nil {
		return nil, err
	}
	return &model.PageResult[model.MessageModel]{
		List:  list,
		Count: count,
	}, nil
}
