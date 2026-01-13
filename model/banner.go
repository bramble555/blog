package model

import (
	"os"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model/ctype"
	"gorm.io/gorm"
)

type BannerModel struct {
	MODEL
	Hash      string           `json:"hash"`                        // 图片的hash值，用于判断重复图片
	Name      string           `gorm:"size:38" json:"name"`         // 图片名称
	ImageType ctype.BannerType `gorm:"default:1" json:"image_type"` // 图片类型，本地还是网上的,1 是本地
}

func (b *BannerModel) BeforeDelete(tx *gorm.DB) (err error) {
	if b.ImageType == ctype.Local {
		// 如果是本地图片,删除本地存储图片
		err = os.Remove(global.Config.Upload.Path + "/" + b.Name)
		if err != nil {
			global.Log.Errorf("BeforeDelete err:%s\n", err.Error())
			return err
		}
	}
	return nil
}

type FileUploadResponse struct {
	FileName  string `json:"file_name"`
	IsSuccess bool   `json:"is_success"`
	Msg       string `json:"msg"`
}
type ResponseBanner struct {
	SN   int64  `json:"sn,string"`
	Name string `gorm:"size:38" json:"name"`
}

var WhiteImageExtList = map[string]struct{}{
	"jpg":  {},
	"jpeg": {},
	"png":  {},
	"ico":  {},
	"tiff": {},
	"gif":  {},
	"svg":  {},
	"webg": {},
}
