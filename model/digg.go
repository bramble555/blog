package model

import "time"

// UserDiggModel 用户文章点赞表
type UserDiggModel struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserSN     int64     `gorm:"index:idx_user_article,unique;not null" json:"user_sn"`
	ArticleSN  int64     `gorm:"index:idx_user_article,unique;not null" json:"article_sn"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
}

func (UserDiggModel) TableName() string {
	return "user_digg_models"
}

// UserCommentDiggModel 用户评论点赞表
type UserCommentDiggModel struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserSN     int64     `gorm:"index:idx_user_comment,unique;not null" json:"user_sn"`
	CommentSN  int64     `gorm:"index:idx_user_comment,unique;not null" json:"comment_sn"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
}

func (UserCommentDiggModel) TableName() string {
	return "user_comment_digg_models"
}
