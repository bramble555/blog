package advert

import (
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func UpdateAdvertShow(p *model.ParamUpdateAdvertShow) (string, error) {
	err := global.DB.Table("advert_models").Where("sn = ?", p.SN).
		Update("is_show", p.IsShow).
		Error
	if err != nil {
		return "", err
	}
	return code.StrUpdateSucceed, nil
}
