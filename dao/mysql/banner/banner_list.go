package banner

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/image"
)

func GetBannerList(il *image.ParamImageList) ([]model.BannerModel, error) {
	var bml []model.BannerModel
	offset := (il.Page - 1) * il.Size
	// SELECT * FROM `banner_models` ORDER BY create_time DESC LIMIT X OFFSET Y
	// 默认降序排列
	err := global.DB.Table("banner_models").Order(il.Order).Limit(il.Size).
		Offset(offset).Find(&bml).Error
	if err != nil {
		global.Log.Errorf("global.DB.Table(banner_models) err:%s\n", err.Error())
		return nil, err
	}
	return bml, nil
}
func DeleteBannerList(pdl *model.ParamDeleteList) (string, error) {
	var bannerList []model.BannerModel
	// SELECT * FROM banner_list WHERE id IN (id1, id2, id3, ...);
	count := global.DB.Find(&bannerList, pdl.IDList).RowsAffected
	if count == 0 {
		global.Log.Errorf("图片不存在")
	}
	// DELETE FROM banner_list WHERE id IN (id1, id2, id3, ...);
	global.DB.Delete(&bannerList)
	return fmt.Sprintf("共删除%d图片", count), nil
}
