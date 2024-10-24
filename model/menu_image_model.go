package model

// 菜单和背景图的连接表，方便排序
type MenuBannerModel struct {
	MenuID   uint `json:"menu_id"`
	BannerID uint `json:"banner_id"`
	Sort     int  `gorm:"size:10" json:"sort"`
}
