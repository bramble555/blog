package banner

import (
	"fmt"

	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"gorm.io/gorm"
)

func GetBannerList(pl *model.ParamList) (*[]model.BannerModel, error) {
	bml, err := mysql.GetTableList[model.BannerModel]("banner_models", pl, "")
	if err != nil {
		return nil, err
	}
	return &bml, nil
}
func GetBannerDetail() (*[]model.ResponseBanner, error) {
	var bd []model.ResponseBanner
	err := global.DB.Table("banner_models").Select("id, name").Scan(&bd).Error
	if err != nil {
		return nil, err
	}
	return &bd, nil
}
func GetBannerByID(id *uint) (*model.ResponseBanner, error) {
	var bd model.ResponseBanner
	err := global.DB.Table("banner_models").Select("id, name").Where("id = ?", id).First(&bd).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.ErrorIDExit
		}
		global.Log.Errorf("banner GetBannerDetail err:%s\n", err.Error())
		return nil, err
	}
	return &bd, nil
}
func DeleteBannerList(pdl *model.ParamDeleteList) (string, error) {
	// 手动开启事务
	t := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			t.Rollback()
			global.Log.Errorf("发生错误: %v\n", r)
		}
	}()

	// 删除 banner_models
	var bd []model.BannerModel
	if err := t.Where("id IN ?", pdl.IDList).Find(&bd).Error; err != nil {
		t.Rollback()
		global.Log.Errorf("查询 banner_models 时出错:%s\n", err.Error())
		return "查询 banner_models 时出错", err
	}

	resultB := t.Delete(&bd)
	if resultB.Error != nil {
		t.Rollback()
		global.Log.Errorf("删除 banner_models 时出错:%v\n", resultB.Error)
		return "删除 banner_models 时出错", code.ErrorIDNotExit
	}

	// 删除 menu_models
	var mm []model.MenuModel
	resultM := t.Where("banner_id IN ?", pdl.IDList).Find(&mm)
	if resultM.Error != nil {
		t.Rollback()
		global.Log.Errorf("查询 menu_models 时出错:%s\n", resultM.Error)
		return "查询 menu_models 时出错", resultM.Error
	}

	if resultM.RowsAffected > 0 {
		// 使用条件直接删除，而不是删除切片
		if err := t.Where("banner_id IN ?", pdl.IDList).Delete(&model.MenuModel{}).Error; err != nil {
			t.Rollback()
			global.Log.Errorf("删除 menu_models 时出错:%v\n", err)
			return "删除 menu_models 时出错", err
		}
	}

	t.Commit()
	return fmt.Sprintf("banner_models 共删除 %d 行, menu_models 共删除 %d 行",
		resultB.RowsAffected, resultM.RowsAffected), nil
}
