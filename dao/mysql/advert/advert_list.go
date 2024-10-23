package advert

import (
	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func GetAdvertList(pl *model.ParamList, isShow bool) ([]model.AdvertModel, error) {
	ad, err := mysql.GetTableList[model.AdvertModel]("advert_models", pl, "")
	if err != nil {
		global.Log.Errorf("advert global.DB.Table(advert_models) err:%s\n", err.Error())
		return nil, err
	}
	// 如果是 admin，直接返回所有广告
	if isShow {
		return ad, nil
	}
	// 筛选出 IsShow 为 true 的广告
	res := make([]model.AdvertModel, 0, len(ad)) // 初始化切片，预分配容量
	for _, advert := range ad {
		if advert.IsShow {
			res = append(res, advert)
		}
	}
	return res, nil
}
