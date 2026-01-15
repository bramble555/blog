package banner

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

// CheckBannerNotExists 检查图片是否存在数据库，如果存在，返回false，不存在返回true
func CheckBannerNotExists(hash string) bool {
	var count int64
	global.DB.Table("banner_models").Where("hash=?", hash).Count(&count)
	return count == 0
}
func UploadBanners(hash string, fileName string) error {
	// 尝试将图片信息插入数据库
	err := global.DB.Create(&model.BannerModel{
		Hash: hash,
		Name: fileName,
	}).Error
	if err != nil {
		global.Log.Errorf("将图片 %s 插入数据库时发生错误: %v", fileName, err)
		return err
	}
	return nil
}
