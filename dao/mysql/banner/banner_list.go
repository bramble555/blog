package banner

import (
	"fmt"

	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func GetBannerList(pl *model.ParamList) ([]model.BannerModel, error) {
	bml, err := mysql.GetTableList[model.BannerModel]("banner_models", pl, "")
	if err != nil {
		global.Log.Errorf("banner global.DB.Table(banner_models) err:%s\n", err.Error())
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
	return fmt.Sprintf("共删除%d张图片", count), nil
}
