package banner

import (
	"fmt"

	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
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
	err := global.DB.Table("banner_models").Select("sn, name").Scan(&bd).Error
	if err != nil {
		return nil, err
	}
	return &bd, nil
}
func GetBannerBySN(sn *int64) (*model.ResponseBanner, error) {
	var bd model.ResponseBanner
	err := global.DB.Table("banner_models").Select("sn, name").Where("sn = ?", sn).First(&bd).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.ErrorSNNotExit
		}
		global.Log.Errorf("banner GetBannerDetail err:%s\n", err.Error())
		return nil, err
	}
	return &bd, nil
}
func DeleteBannerList(pdl *model.ParamDeleteList) (string, error) {
	// 转换 SNList 为 []int64
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("DeleteBannerList StringSliceToInt64Slice err: %s\n", err.Error())
		return "", err
	}

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
	if err := t.Where("sn IN ?", snList).Find(&bd).Error; err != nil {
		t.Rollback()
		global.Log.Errorf("查询 banner_models 时出错:%s\n", err.Error())
		return "查询 banner_models 时出错", err
	}

	resultB := t.Delete(&bd)
	if resultB.Error != nil {
		t.Rollback()
		global.Log.Errorf("删除 banner_models 时出错:%v\n", resultB.Error)
		return "删除 banner_models 时出错", code.ErrorSNNotExit
	}

	t.Commit()
	return fmt.Sprintf("banner_models 共删除 %d 行", resultB.RowsAffected), nil
}
