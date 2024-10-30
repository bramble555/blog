package logic

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"github.com/bramble555/blog/pkg"
)

// 用户名或者邮箱登录
func UsernameLogin(peu *model.ParamEmailUser) (string, error) {
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
func GetUserList(role uint, pl *model.ParamList) (*[]model.UserModel, error) {
	udl, err := user.GetUserList(pl)
	if err != nil {
		return nil, err
	}

	// 手机号和邮箱脱敏处理
	for i := range *udl {
		(*udl)[i].Email = pkg.DesensitizeEmail((*udl)[i].Email)
		(*udl)[i].Phone = pkg.DesensitizePhone((*udl)[i].Phone)
	}

	// 如果是普通用户，username 返回 "****"
	if role == uint(ctype.PermissionUser) {
		for i := range *udl {
			(*udl)[i].Username = "****"
		}
	}

	return udl, nil
}

