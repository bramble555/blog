package model

import "time"

type UserCollectModel struct {
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`

	UserSN    int64 `json:"user_sn,string"`
	ArticleSN int64 `json:"article_sn,string"`
}
