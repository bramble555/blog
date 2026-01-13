package mysql

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

// DeleteTableList 获取列表
func GetTableList[T any](tableName string, pl *model.ParamList, where string, args ...any) ([]T, error) {
	var results []T
	if pl != nil {
		offset := (pl.Page - 1) * pl.Size
		err := global.DB.Table(tableName).
			Where(where, args...).
			Order(pl.Order).
			Limit(pl.Size).
			Offset(offset).
			Find(&results).Error
		return results, err
	}
	// 如果 pl 是 nil, 那就查询所有，不是分页查询
	err := global.DB.Table(tableName).Where(where, args...).Find(&results).Error
	if err != nil {
		global.Log.Errorf("查询时候出错: err:%v\n", err.Error())
		return nil, err
	}
	return results, nil
}

// DeleteTableList 删除列表
func DeleteTableList[T any](tableName string, pdl *model.ParamDeleteList) (string, error) {
	// 转换 SNList 为 []int64
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("DeleteTableList StringSliceToInt64Slice err: %s\n", err.Error())
		return "", err
	}

	var records []T
	// 查找记录
	if err := global.DB.Where("sn IN ?", snList).Find(&records).Error; err != nil {
		return "", err
	}
	// 要启动钩子函数，必须要先查询，然后delete，使用指针类型的 records 进行删除
	result := global.DB.Delete(&records)
	if result.Error != nil {
		global.Log.Errorf("删除记录时出错: %v\n", result.Error)
		return "", fmt.Errorf("要删除的 SN 列表 %v 不存在", snList)

	}
	return fmt.Sprintf("共删除%d条记录", result.RowsAffected), nil
}
