package model

import "time"

type MODEL struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}
type ParamDeleteList struct {
	IDList int64 `json:"id_list" binding:"required"`
}
