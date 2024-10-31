package user

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func UpdateUserRole(puur *model.ParamUpdateUserRole) (string, error) {
	err := global.DB.Table("user_models").Update("role", puur.Role).
		Where("id = ?", puur.UserID).Error
	if err != nil {
		global.Log.Errorf("user UpdateUserRole err:%s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("修改用户%d权限成功", puur.UserID), nil
}

// UpdateUserPwd 负责更新用户密码
func UpdateUserPwd(puup *model.ParamUpdateUserPwd, id uint) (string, error) {
	err := global.DB.Table("user_models").Where("id = ?", id).
		Update("password", puup.Pwd).Error
	if err != nil {
		global.Log.Errorf("user UpdateUserPwd err: %s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("修改用户 %d 密码成功", id), nil
}
