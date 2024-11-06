package model

import "time"

// ArticleTagModel 表示文章和标签之间的关联关系
type ArticleTagModel struct {
	ArticleID  uint      `gorm:"primary_key"`
	TagID      uint      `gorm:"primary_key"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;on_update:CURRENT_TIMESTAMP"`
}
