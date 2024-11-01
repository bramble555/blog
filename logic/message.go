package logic

import (
	errcode "github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/message"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/model"
)

func SendMessage(pm *model.ParamMessage) (string, error) {
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
	return message.SendMessage(pm)
}
