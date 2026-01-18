package tag

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"gorm.io/gorm"
)

func CreateTags(tm *model.TagModel) (string, error) {
	// 判断tag是否存在
	err := global.DB.Where("title = ?", tm.Title).First(&tm).Error

	if err == nil {
		// 找到记录 → 标签已存在
		return "", code.ErrorTagExist
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 非 ErrRecordNotFound → 数据库错误
		return "", err
	}
	// 创建tag
	err = global.DB.Create(&model.TagModel{
		Title: tm.Title,
	}).Error
	if err != nil {
		global.Log.Errorf("mysql CreateTag err:%s\n", err.Error())
		return "", err
	}
	return code.StrCreateSucceed, nil
}
