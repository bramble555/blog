package logic

import (
	errcode "github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/message"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/model"
)

func SendMessage(pm *model.ParamMessage) (string, error) {
	// 先判断发送方和接收方是否存在
	ok, err := user.CheckUserExistBySN(int64(pm.SendUserSN))
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errcode.ErrorUserNotExit
	}
	ok, err = user.CheckUserExistBySN(int64(pm.RevUserSN))
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errcode.ErrorUserNotExit
	}
	sud, err := user.GetUserDetailBySN(int64(pm.SendUserSN))
	if err != nil {
		return "", err
	}
	rud, err := user.GetUserDetailBySN(int64(pm.RevUserSN))
	if err != nil {
		return "", err
	}
	udl := make([]*model.UserDetail, 2)
	udl[0] = sud
	udl[1] = rud
	return message.SendMessage(pm)
}
func MessageListAll(pl *model.ParamList) ([]model.RespondMessage, error) {
	return message.MessageListAll(pl)

}
func MessageList(sn int64) ([]model.RespondMessage, error) {
	return message.MessageList(sn)
}
func MessageRecord(mySN, recordSN int64) ([]model.MessageModel, error) {
	return message.MessageRecord(mySN, recordSN)
}
