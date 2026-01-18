package model

import (
	"time"

	"github.com/bramble555/blog/global"
	"gorm.io/gorm"
)

type LoginModel struct {
	SN         int64     `json:"sn"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	Username   string    `json:"username"`
}

func (l *LoginModel) BeforeCreate(tx *gorm.DB) error {
	// 防御性赋值
	if l.SN == 0 {
		l.SN = global.Snowflake.GetID()
	}
	return nil
}
func (LoginModel) TableName() string {
	return "login_models"
}
