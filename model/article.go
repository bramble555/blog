package model

import "github.com/bramble555/blog/model/ctype"

type ArticleModel struct {
	MODEL
	Title         string      `json:"title"`          // 文章标题
	Abstract      string      `json:"abstract"`       // 文章简介
	Content       string      `json:"content"`        // 文章内容
	LookCount     uint        `json:"look_count"`     // 浏览量
	CommentCount  uint        `json:"comment_count"`  // 评论量
	DiggCount     uint        `json:"digg_count"`     // 点赞量
	CollectsCount uint        `json:"collects_count"` // 收藏量
	Category      string      `json:"category"`       // 文章分类
	Source        string      `json:"source"`         // 文章来源
	Link          string      `json:"link"`           // 原文链接
	BannerID      uint        `json:"banner_id"`      // 文章封面ID
	Tags          ctype.Array `json:"tags"`           // 文章标签
	UserID        uint        `json:"user_id"`        // 用户id
}
