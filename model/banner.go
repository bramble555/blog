package model

import (
	"github.com/bramble555/blog/model/ctype"
	"gorm.io/gorm"
)

type BannerModel struct {
	MODEL
	Hash      string           `json:"hash"`                        // 图片的hash值，用于判断重复图片
	Name      string           `gorm:"size:38" json:"name"`         // 图片名称
	ImageType ctype.BannerType `gorm:"default:1" json:"image_type"` // 图片类型，本地还是网上的,1 是本地
	Path      string           `json:"path" gorm:"-"`
}

func (BannerModel) TableName() string {
	return "banner_models"
}

func (b *BannerModel) AfterFind(tx *gorm.DB) (err error) {
	if b.ImageType == ctype.Local {
		b.Path = "/uploads/file/" + b.Name
	}
	return
}

type FileUploadResponse struct {
	FileName string `json:"file_name"`
	Msg      string `json:"msg"`
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
