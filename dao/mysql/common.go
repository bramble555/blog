package mysql

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func GetTableList[T any](tableName string, pl *model.ParamList, where string, args ...interface{}) ([]T, error) {
	var results []T
	offset := (pl.Page - 1) * pl.Size
	err := global.DB.Table(tableName).
		Where(where, args...).
		Order(pl.Order).
		Limit(pl.Size).
		Offset(offset).
		Find(&results).Error
	
	return results, err
}

