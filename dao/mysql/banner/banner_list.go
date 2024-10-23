package banner

import (
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
	return mysql.DeleteTableList[model.BannerModel]("banner_models", pdl)
}
