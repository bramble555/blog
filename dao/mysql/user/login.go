package user

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

// CheckUserExistByName 检查用户是否存在
func CheckUserExistByName(name string) (bool, error) {
	var count int64
	err := global.DB.Table("user_models").Where("username = ?", name).Count(&count).Error
	if err != nil {
		global.Log.Errorf("user CheckUserExistByName err:%s\n", err.Error())
		return false, err
	}
	return count == 1, nil
}
func CheckUserExistByID(id uint) (bool, error) {
	var count int64
	err := global.DB.Table("user_models").Where("id = ?", id).Count(&count).Error
	if err != nil {
		global.Log.Errorf("user CheckUserExistByID err:%s\n", err.Error())
		return false, err
	}
	return count == 1, nil
}

// CheckPwdExistByID 传入 ID，检查密码是否正确
func CheckPwdExistByID(id uint, pwd string) (bool, error) {
	var encryPassword string
	err := global.DB.Table("user_models").Where("id = ?", id).
		Select("password").Scan(&encryPassword).Error
	if err != nil {
		global.Log.Errorf("user QueryPassword err: %v\n", err)
		return false, err
	}

	// 如果密码为空，用户不存在
	if encryPassword == "" {
		return false, errors.New("用户不存在")
	}

	// 比较密码
	err = pkg.ComparePasswords(encryPassword, pwd)
	if err != nil {
		return false, errors.New("密码不正确")
	}
	return true, nil
}

// QueryPasswordByUsername 传入username 和 密码 检验密码是否正确，实现登录功能
func QueryPasswordByUsername(peu *model.ParamEmailUser) (bool, error) {
	var encryPassword string
	err := global.DB.Table("user_models").Where("username = ?", peu.Username).
		Select("password").Scan(&encryPassword).Error // 使用 Scan 将结果绑定到 password
	if err != nil {
		global.Log.Errorf("user QueryPassword err: %v\n", err)
		return false, err
	}
	err = pkg.ComparePasswords(encryPassword, peu.Password)
	if err != nil {
		global.Log.Errorf("user pkg.ComparePassword serr: %v\n", err)
		return false, code.ErrorPasswordWrong
	}
	return true, nil
}
func GetToken(peu *model.ParamEmailUser) (string, error) {
	type paramUserDetail struct {
		ID       uint // 改为大写 否则不能 Scan 到
		Username string
		Role     uint // 改为大写
	}
	var udd paramUserDetail

	err := global.DB.Table("user_models").Where("username = ?", peu.Username).
		Select("id,username,role").Scan(&udd).Error
	if err != nil {
		global.Log.Errorf("user GetToken select err:%s\n", err.Error())
		return "", err
	}
	token, err := pkg.GenToken(udd.ID, udd.Role, udd.Username)
	if err != nil {
		global.Log.Errorf("pkg GetToken err:%s\n", err.Error())
		return "", err
	}
	return token, nil
}
