package advert

import (
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func CreateAdvert(ad *model.AdvertModel) (string, error) {
	// 在创建之前先检查是否存在
	existingAdvert := model.AdvertModel{}
	existingCount := global.DB.Take(&existingAdvert, "title = ?", ad.Title).RowsAffected
	if existingCount == 1 {
		global.Log.Errorf("mysql global.DB.Create(&model.AdvertModel error\n")
		return "", code.ErrorTitleExist
	}
	global.DB.Table("advert_models").Where("title = ?", ad.Title)
	// 不存在就创建
	err := global.DB.Create(&model.AdvertModel{
		Title:  ad.Title,
		Href:   ad.Href,
		Images: ad.Images,
		IsShow: ad.IsShow,
	}).Error
	if err != nil {
		global.Log.Errorf("mysql global.DB.Create&model.AdvertModel err:%s\n", err.Error())
		return "", err
	}
	return code.StrCreateSucceed, nil
}
