package user

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

func UpdateUserRole(puur *model.ParamUpdateUserRole) (string, error) {
	global.Log.Debugf("id:%d", puur.UserID)
	err := global.DB.Table("user_models").Where("id = ?", puur.UserID).Update("role", puur.Role).Error
	if err != nil {
		global.Log.Errorf("user UpdateUserRole err:%s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("修改用户%d权限成功", puur.UserID), nil
}

// UpdateUserPwd 负责更新用户密码
func UpdateUserPwd(puup *model.ParamUpdateUserPwd, id uint) (string, error) {
	// 先把先密码加密
	pwd, err := pkg.HashPassword(puup.Pwd)
	if err != nil {
		global.Log.Errorf("HashPassword err:%s\n", err.Error())
		return "", err
	}
	err = global.DB.Table("user_models").Where("id = ?", id).
		Update("password", pwd).Error
	if err != nil {
		global.Log.Errorf("user UpdateUserPwd err: %s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("修改用户 %d 密码成功", id), nil
}
