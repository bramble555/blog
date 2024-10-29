package logic

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/model"
)

// EmailLogin 邮箱登录
func EmailLogin(peu *model.ParamEmailUser) (string, error) {
	// 判断用户名是否存在
	ok, err := user.CheckUserExist(peu.Username)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("用户名不存在")
	}
	// 判断密码是否错误
	ok, err = user.QueryPassword(peu)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("密码错误")
	}
	return user.GetToken(peu)
}
