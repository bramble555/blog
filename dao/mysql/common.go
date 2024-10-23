package mysql

import (
	"fmt"

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
func DeleteTableList[T any](tableName string, pdl *model.ParamDeleteList) (string, error) {
	var records []T
	if len(pdl.IDList) == 0 {
		return "", fmt.Errorf("没有提供 ID")
	}
	result := global.DB.Where("id IN ?", pdl.IDList).Delete(records)
	if result.Error != nil {
		global.Log.Errorf("删除记录时出错: %v\n", result.Error)
		return "", result.Error
	}
	return fmt.Sprintf("共删除%d条记录", result.RowsAffected), nil
}
