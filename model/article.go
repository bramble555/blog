package model

import (
	"github.com/bramble555/blog/model/ctype"
)

type ArticleModel struct {
	MODEL
	Title         string      `json:"title"`          // 文章标题
	Abstract      string      `json:"abstract"`       // 文章简介
	Content       string      `json:"content"`        // 文章内容
	LookCount     int64       `json:"look_count"`     // 浏览量
	CommentCount  int64       `json:"comment_count"`  // 评论量
	DiggCount     int64       `json:"digg_count"`     // 点赞量
	CollectsCount int64       `json:"collects_count"` // 收藏量
	Category      string      `json:"category"`       // 文章分类
	Source        string      `json:"source"`         // 文章来源
	Link          string      `json:"link"`           // 原文链接
	Tags          ctype.Array `json:"tags"`           // 文章标签

	BannerSN  int64  `json:"banner_sn,string"` // 文章封面 SN
	BannerUrl string `json:"banner_url"`

	UserSN     int64  `json:"user_sn,string"` // 用户 SN
	Username   string `json:"username"`
	UserAvatar string `json:"user_avatar"`

	IsCollect bool `gorm:"-" json:"is_collect"` // 是否收藏，不入库
}
type ParamArticle struct {
	Title    string      `json:"title" binding:"required"`   // 文章标题
	Abstract string      `json:"abstract"`                   // 文章简介
	Content  string      `json:"content" binding:"required"` // 文章内容
	Category string      `json:"category"`                   // 文章分类
	Source   string      `json:"source"`                     // 文章来源
	Link     string      `json:"link"`                       // 原文链接
	Tags     ctype.Array `json:"tags"`                       // 文章标签
	BannerSN int64       `json:"banner_sn,string"`           // 文章封面 SN
}
type ResponseArticle struct {
	SN            int64       `gorm:"column:sn" json:"sn,string"`
	CreateTime    string      `json:"create_time"`
	UpdateTime    string      `json:"update_time"`
	Title         string      `json:"title"`          // 文章标题
	Abstract      string      `json:"abstract"`       // 文章简介
	Content       string      `json:"content"`        // 文章内容
	LookCount     int64       `json:"look_count"`     // 浏览量
	CommentCount  int64       `json:"comment_count"`  // 评论量
	DiggCount     int64       `json:"digg_count"`     // 点赞量
	CollectsCount int64       `json:"collects_count"` // 收藏量
	Category      string      `json:"category"`       // 文章分类
	Tags          ctype.Array `json:"tags"`           // 文章标签

	BannerSN  int64  `json:"banner_sn"` // 文章封面 SN
	BannerUrl string `json:"banner_url"`

	UserSN     int64  `json:"user_sn"` // 用户 SN
	Username   string `json:"username"`
	UserAvatar string `json:"user_avatar"`

	IsCollect bool `gorm:"-" json:"is_collect"` // 是否收藏，不入库
}

// func (ra *ResponseArticle) ParseTags() ([]string, error) {
// 	var tags []string
// 	if err := json.Unmarshal(ra.Tags, &tags); err != nil {
// 		// 如果是字符串而不是数组，尝试将其解析为单一字符串
// 		var tag string
// 		if err := json.Unmarshal(ra.Tags, &tag); err != nil {
// 			return nil, err
// 		}
// 		tags = []string{tag}
// 	}
// 	return tags, nil
// }

type CalendarCount struct {
	Date  string `json:"data"`
	Count int    `json:"count"`
}
type ParamArticleQuery struct {
	ParamList
	Title   string `json:"title" form:"title"`
	Tags    string `json:"tags" form:"tags"`
	Content string `json:"content" form:"content"`
}
