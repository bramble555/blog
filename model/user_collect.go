package model

import "time"

type UserCollectModel struct {
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`

	UserID    uint `json:"user_id"`
	ArticleID uint `json:"article_id"`
}
