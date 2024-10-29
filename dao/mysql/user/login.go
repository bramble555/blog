package user

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

// CheckUserExist 检查用户是否存在
func CheckUserExist(name string) (bool, error) {
	var count int64
	err := global.DB.Table("user_models").Where("username = ?", name).Count(&count).Error
	if err != nil {
		global.Log.Errorf("user CheckUserExist err:%s\n", err.Error())
		return false, err
	}
	return count == 1, nil
}

// QueryPassword 检验密码是否正确
func QueryPassword(peu *model.ParamEmailUser) (bool, error) {
	var encryPassword string
	err := global.DB.Table("user_models").Where("username = ?", peu.Username).
		Select("password").Scan(&encryPassword).Error // 使用 Scan 将结果绑定到 password
	if err != nil {
		global.Log.Errorf("user QueryPassword err: %v\n", err)
		return false, err
	}
	err = pkg.ComparePasswords(encryPassword, peu.Password)
	if err != nil {
		return false, err
	}
	return true, nil
}
func GetToken(peu *model.ParamEmailUser) (string, error) {
	type userModelDetail struct {
		id   int
		role int
	}
	var udd userModelDetail
	err := global.DB.Table("user_models").Where("username = ?", peu.Username).
		Select("id, role").Scan(&udd).Error
	if err != nil {
		global.Log.Errorf("user GetToken select err:%s\n", err.Error())
		return "", err
	}
	token, err := pkg.GenToken(udd.id, udd.role)
	if err != nil {
		global.Log.Errorf("pkg GetToken err:%s\n", err.Error())
		return "", err
	}
	return token, nil
}
