package tag

import (
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func CreateTags(ad *model.TagModel) (string, error) {
	err := global.DB.Create(&model.TagModel{
		Title: ad.Title,
	}).Error
	if err != nil {
		global.Log.Errorf("mysql CreateAdvert err:%s\n", err.Error())
		return "", err
	}
	return code.CreateSucceed, nil
}
