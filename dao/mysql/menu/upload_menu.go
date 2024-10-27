package menu

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func UploadMenu(mm *model.MenuModel) (string, error) {
	err := global.DB.Create(mm).Error
	if err != nil {
		global.Log.Errorf("menu create err:%s\n", err.Error())
		return "创建菜单失败", err
	}
	return "创建菜单成功", nil
}
