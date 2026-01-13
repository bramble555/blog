package user

import (
	"fmt"

	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

func UpdateUserRole(puur *model.ParamUpdateUserRole) (string, error) {
	global.Log.Debugf("sn:%d", puur.UserSN)
	err := global.DB.Table("user_models").Where("sn = ?", puur.UserSN).Update("role", puur.Role).Error
	if err != nil {
		global.Log.Errorf("user UpdateUserRole err:%s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("修改用户%d权限成功", puur.UserSN), nil
}

// UpdateUserPwd 负责更新用户密码
func UpdateUserPwd(puup *model.ParamUpdateUserPwd, sn int64) (string, error) {
	// 先把先密码加密
	pwd, err := pkg.HashPassword(puup.Pwd)
	if err != nil {
		global.Log.Errorf("HashPassword err:%s\n", err.Error())
		return "", err
	}
	err = global.DB.Table("user_models").Where("sn = ?", sn).
		Update("password", pwd).Error
	if err != nil {
		global.Log.Errorf("user UpdateUserPwd err: %s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("修改用户 %d 密码成功", sn), nil
}

// DeleteUserList 删除用户列表
func DeleteUserList(pdl *model.ParamDeleteList) (string, error) {
	// 其实还要删除许多关联的表，后面再删除
	// to do
	return mysql.DeleteTableList[model.UserModel]("user_models", pdl)
}
