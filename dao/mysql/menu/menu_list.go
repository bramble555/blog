package menu

import (
	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func GetMenuList() (*[]model.ResponseMenuBanner, error) {
	// 获取菜单列表全部
	mm, err := mysql.GetTableList[model.MenuModel]("menu_models", nil, "")
	if err != nil {
		global.Log.Errorf("menu GetMenuList err: %s\n", err.Error())
		return nil, err
	}
	// 获取 Banner ID 列表
	idList := make([]*uint, len(mm))
	for i, menu := range mm {
		idList[i] = menu.BannerID
	}
	// 获取 Banner 数据
	var banners []model.ResponseBanner
	if len(idList) > 0 {
		err = global.DB.Table("banner_models").
			Where("id IN (?)", idList).
			Select("id, name").Find(&banners).Error
		if err != nil {
			global.Log.Errorf("banner GetMenuList err: %s\n", err.Error())
			return nil, err
		}
	}
	n := len(mm)
	// 组装 ResponseMenuBanner 数据
	response := make([]model.ResponseMenuBanner, n)
	for i := 0; i < n; i++ {
		response[i].MenuModel = &mm[i]
		response[i].Name = banners[i].Name
	}
	return &response, nil
}
func DeleteMenuList(pdl *model.ParamDeleteList) (string, error) {
	return mysql.DeleteTableList[model.ParamDeleteList]("menu_models", pdl)
}
