package advert

import (
	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/model"
)

func GetAdvertList(pl *model.ParamList, isShow bool) ([]model.AdvertModel, error) {
	var condition string
	var args []any
	// isShow = true 返回 数据库中 isShow = true or false
	// isShow = false 返回数据库中 isShow = true
	if !isShow {
		condition = "is_show = ?"
		args = append(args, true) // 添加过滤条件的值
	}
	// 调用 GetTableList 时，只有在 condition 不为空时才传递条件和参数
	ad, err := mysql.GetTableList[model.AdvertModel]("advert_models", pl, condition, args...)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func DeleteAdvertList(pdl *model.ParamDeleteList) (string, error) {
	return mysql.DeleteTableList[model.AdvertModel]("advert_models", pdl)
}
