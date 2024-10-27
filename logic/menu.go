package logic

import (
	"github.com/bramble555/blog/dao/mysql/banner"
	"github.com/bramble555/blog/dao/mysql/menu"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func UploadMenu(mm *model.MenuModel) (string, error) {
	// 查询用户传递的 banner_id 是否存在
	if mm.BannerID != nil {
		_, err := banner.GetBannerByID(mm.BannerID)
		if err != nil {
			global.Log.Errorf("banner GetBannerByID err:%s\n", err.Error())
			return "", err
		}
	}
	global.Log.Debugf("banner_id:%d", mm.BannerID)
	return menu.UploadMenu(mm)
}
func GetMenuList() (*[]model.ResponseMenuBanner, error) {
	return menu.GetMenuList()
}
func UpdateMenu(id uint, mm *model.MenuModel) (string, error) {
	return menu.UpdateMenu(id, mm)
}
func DeleteMenuList(pdl *model.ParamDeleteList) (string, error) {
	return menu.DeleteMenuList(pdl)
}
