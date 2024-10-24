package model

// 菜单表 菜单的的路径可以是 /path 也可以是路由别名
type MenuModel struct {
	MODEL
	Title        string `gorm:"size:32" json:"title"`
	Path         string `gorm:"size:32" json:"path"`
	Slogan       string `gorm:"size:64" json:"slogan"` // 口号，标语
	Abstract     string `json:"abstract"`              // 简介
	AbstractTime uint   `json:"abstract_time"`         // 简介的切换时间 单位为秒
	BannerID     *uint  `json:"banner_id"`             // 指针类型，允许为 nil
	BannerTime   uint   `json:"banner_time"`           // 菜单图片的切换时间为多少秒 0表示不切换
	Sort         uint   `gorm:"size:5" json:"sort"`    // 菜单的顺序 0 1 2
}
