package banner

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

// CheckBannerNotExists 检查图片是否存在数据库，如果存在，返回false，不存在返回true
func CheckBannerNotExists(byteData []byte) bool {
	hash := pkg.MD5(byteData)
	var sn int64
	// select sn from banner where hash = hash
	rows := global.DB.Table("banner_models").Select("sn").Where("hash=?", hash).Scan(&sn).RowsAffected
	return rows != 1
}
func UploadBanners(byteData []byte, fileName string) error {
	// 生成文件的 MD5 哈希值
	hash := pkg.MD5(byteData)

	// 尝试将图片信息插入数据库
	result := global.DB.Create(&model.BannerModel{
		Hash: hash,
		Name: fileName,
	})

	// 检查是否有错误
	if result.Error != nil {
		global.Log.Errorf("将图片 %s 插入数据库时发生错误: %v", fileName, result.Error)
		return result.Error
	}

	// 检查是否影响了 1 行
	if result.RowsAffected != 1 {
		global.Log.Errorf("将图片 %s 插入数据库时未成功插入 (影响行数: %d)", fileName, result.RowsAffected)
		return fmt.Errorf("插入图片 %s 时影响行数不正确", fileName)
	}
	return nil
}
