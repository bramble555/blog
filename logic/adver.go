package logic

import (
	"github.com/bramble555/blog/dao/mysql/advert"
	"github.com/bramble555/blog/model"
)

func CreateAdvert(ad *model.AdvertModel) (string, error) {
	return advert.CreateAdvert(ad)
}
