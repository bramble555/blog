package logic

import (
	errcode "github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/message"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/model"
)

func SendMessage(pm *model.ParamMessage) (string, error) {
	// 先判断发送方和接收方是否存在
	ok, err := user.CheckUserExistByID(pm.SendUserID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errcode.ErrorUserNotExit
	}
	ok, err = user.CheckUserExistByID(pm.RevUserID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errcode.ErrorUserNotExit
	}
	sud, err := user.GetUserDetail(pm.SendUserID)
	if err != nil {
		return "", err
	}
	rud, err := user.GetUserDetail(pm.RevUserID)
	udl := make([]*model.UserDetail, 2)
	udl[0] = sud
	udl[1] = rud
	return message.SendMessage(pm, udl)
}
func MessageListAll(pl *model.ParamList) ([]model.RespondMessage, error) {
	return message.MessageListAll(pl)

}
func MessageList(id uint) ([]model.RespondMessage, error) {
	return message.MessageList(id)
}
func MessageRecord(myID, recordID uint) ([]model.MessageModel, error) {
	return message.MessageRecord(myID, recordID)
}
