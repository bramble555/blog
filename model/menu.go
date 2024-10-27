package model

import "github.com/bramble555/blog/model/ctype"

// 菜单表 菜单的的路径可以是 /path 也可以是路由别名
type MenuModel struct {
	*MODEL
	Title        string      `gorm:"size:32" json:"title" binding:"required"`
	Path         string      `gorm:"size:32" json:"path" binding:"required"`
	Slogan       string      `gorm:"size:64" json:"slogan"`      // 口号，标语
	Abstract     ctype.Array `json:"abstract"`                   // 简介
	AbstractTime int         `json:"abstract_time"`              // 切换的时间，默认 0
	BannerID     *uint       `json:"banner_id,string,omitempty"` // 不传的时候说明不需要图片
	Sort         uint        `gorm:"size:5" json:"sort"`         // 菜单的顺序 0(默认值) 1 2
}

type ResponseMenuBanner struct {
	*MenuModel
	Name string
}
