package advert

import (
	"errors"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func CreateAdvert(ad *model.AdvertModel) (string, error) {
	// 在创建之前先检查是否存在
	existingAdvert := model.AdvertModel{}
	existingCount := global.DB.Take(&existingAdvert, "title = ?", ad.Title).RowsAffected
	if existingCount == 1 {
		global.Log.Errorf("mysql global.DB.Create(&model.AdvertModel error")
		return "", errors.New("广告已存在")
	}
	global.DB.Table("advert_models").Where("title = ?", ad.Title)
	// 不存在就创建
	count := global.DB.Create(&model.AdvertModel{
		Title:  ad.Title,
		Href:   ad.Href,
		Images: ad.Images,
		IsShow: ad.IsShow,
	}).RowsAffected
	if count != 1 {
		global.Log.Errorf("mysql global.DB.Create(&model.AdvertModel error")
		return "", errors.New("mysql 创建广告失败")
	}
	return "创建广告成功", nil
}

