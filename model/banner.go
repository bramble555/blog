package model

import (
	"os"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model/ctype"
	"gorm.io/gorm"
)

type BannerModel struct {
	MODEL
	Hash      string          `json:"hash"`                        // 图片的hash值，用于判断重复图片
	Name      string          `gorm:"size:38" json:"name"`         // 图片名称
	ImageType ctype.ImageType `gorm:"default:1" json:"image_type"` // 图片类型，本地还是网上的
}

func (b *BannerModel) BeforeDelete(tx *gorm.DB) (err error) {
	if b.ImageType == ctype.Local {
		//本地图片，删除，还要删除本地存储
		err = os.Remove(global.Config.Upload.Path + b.Name)
		if err != nil {
			global.Log.Errorf("err:%s\n", err.Error())
			return err
		}
	}
	return nil
}
