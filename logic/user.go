package logic

import (
	"errors"
	"time"

	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"github.com/bramble555/blog/pkg"
)

// 用户名或者邮箱登录
func UsernameLogin(peu *model.ParamEmailUser) (string, error) {
	// 判断用户名是否存在
	ok, err := user.CheckUserExistByName(peu.Username)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("用户名不存在")
	}
	// 判断密码是否错误
	ok, err = user.QueryPasswordByUsername(peu)
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
func UpdateUserRole(puur *model.ParamUpdateUserRole) (string, error) {
	ok, err := user.CheckUserExistByID(puur.UserID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "用户ID不存在", errors.New("用户ID不存在")
	}
	return user.UpdateUserRole(puur)
}
func UpdateUserPwd(puup *model.ParamUpdateUserPwd, id uint) (string, error) {
	ok, err := user.CheckPwdExistByID(id, puup.OldPwd)
	if err != nil || !ok {
		return "你输入的密码有误,请重新尝试", err
	}
	var data string
	data, err = user.UpdateUserPwd(puup, id)
	if err != nil {
		return "服务器繁忙", err
	}
	return data, nil
}

// Logout 针对注销的操作
func Logout(token string, diff time.Duration) error {
	return redis.Logout(token, diff)
}
