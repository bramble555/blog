package user

import (
	"errors"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

const defaultAvatorPath = "./upload/avator/default.png"

func CreateUser(role int64, username, password string) error {
	// 先检验 username 是否存在
	var count int64
	global.DB.Table("user_models").Where("username = ?", username).Count(&count)
	if count != 0 {
		return errors.New("username 重复了")
	}
	// 把密码加密存储
	ps, err := pkg.HashPassword(password)
	if err != nil {
		global.Log.Errorf("HashPassword err:%s\n", err.Error())
		return err
	}
	// to do 设置一个默认头像
	path := defaultAvatorPath
	user := model.ParamFlagUser{
		SN:       global.Snowflake.GetID(),
		Username: username,
		Password: ps,
		Role:     role,
		Avatar:   path,
	}
	err = global.DB.Table("user_models").Create(&user).Error
	if err != nil {
		global.Log.Errorf("user CreateFlagUser err: %s\n", err.Error())
		return err
	}
	return nil
}
