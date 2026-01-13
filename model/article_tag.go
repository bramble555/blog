package model

import "time"

// ArticleTagModel 表示文章和标签之间的关联关系
type ArticleTagModel struct {
	TagTitle     string    `json:"tag_title"`
	ArticleSN    int64     `json:"article_sn,string"`
	ArticleTitle string    `json:"article_title"`
	CreateTime   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdateTime   time.Time `gorm:"default:CURRENT_TIMESTAMP;on_update:CURRENT_TIMESTAMP"`
}
type ResponseArticleTags struct {
	TagTitle         string    `json:"tag_title"`
	Count            int64     `json:"count"`
	ArticleTitleList string    `json:"article_title_list"`
	CreateTime       time.Time `json:"create_time"`
}
