package logic

import (
	"github.com/bramble555/blog/dao/mysql/advert"
	"github.com/bramble555/blog/model"
)

func CreateAdvert(ad *model.AdvertModel) (string, error) {
	return advert.CreateAdvert(ad)
}
func GetAdvertList(pl *model.ParamList, isShow bool) ([]model.AdvertModel, error) {
	return advert.GetAdvertList(pl, isShow)
}
func DeleteAdvertList(pdl *model.ParamDeleteList) (string, error) {
	return advert.DeleteAdvertList(pdl)
}
