package menu

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func UpdateMenu(id uint, mm *model.MenuModel) (string, error) {
	// 更新数据库中的记录
	result := global.DB.Table("menu_models").Where("id = ?", id).Updates(model.MenuModel{
		Title:        mm.Title,
		Path:         mm.Path,
		Slogan:       mm.Slogan,
		Abstract:     mm.Abstract,
		AbstractTime: mm.AbstractTime,
		BannerID:     mm.BannerID,
		Sort:         mm.Sort,
	})
	// 检查是否有错误
	if result.Error != nil {
		return "更新失败", result.Error
	}

	// 检查是否有行被更新
	if result.RowsAffected == 0 {
		return "更新失败", fmt.Errorf("未找到 ID 为 %d 的记录", id)
	}
	return "更新成功", nil
}
