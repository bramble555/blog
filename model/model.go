package model

import "time"

type MODEL struct {
	ID         uint      `gorm:"primaryKey" json:"id,string"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}
type ParamDeleteList struct {
	IDList []uint `json:"id_list" binding:"required"`
}
type ParamList struct {
	Page  int    `json:"page" form:"page"`
	Size  int    `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}
type ParamID struct {
	ID uint `json:"id,string" binding:"required"`
}

// 默认按照创建时间降序排序
const OrderByTime = "create_time DESC, id DESC"
