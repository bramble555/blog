package model

import (
	"github.com/bramble555/blog/model/ctype"
)

type ArticleModel struct {
	MODEL
	Title         string            `json:"title"`          // 文章标题
	Abstract      string            `json:"abstract"`       // 文章简介
	Content       string            `json:"content"`        // 文章内容
	LookCount     int64             `json:"look_count"`     // 浏览量
	CommentCount  int64             `json:"comment_count"`  // 评论量
	DiggCount     int64             `json:"digg_count"`     // 点赞量
	CollectsCount int64             `json:"collects_count"` // 收藏量
	Tags          ctype.ArrayString `json:"tags"`           // 文章标签

	BannerSN  int64  `json:"banner_sn,string"` // 文章封面 SN
	BannerUrl string `json:"banner_url"`

	UserSN     int64  `json:"user_sn,string"` // 用户 SN
	Username   string `json:"username"`
	UserAvatar string `json:"user_avatar"`

	ParsedContent string `gorm:"-" json:"parsed_content"` // 解析后的 HTML 内容

	IsCollect bool `gorm:"-" json:"is_collect"` // 是否收藏，不入库
	IsDigg    bool `gorm:"-" json:"is_digg"`    // 是否点赞，不入库
}

func (ArticleModel) TableName() string {
	return "article_models"
}

type ParamArticle struct {
	Title    string `json:"title" binding:"required"`   // 文章标题
	Abstract string `json:"abstract"`                   // 文章简介
	Content  string `json:"content" binding:"required"` // 文章内容
	Tags     string `json:"tags"`                       // 文章标签，前端传入逗号分隔的字符串
	BannerSN int64  `json:"banner_sn,string"`           // 文章封面 SN
}
type ResponseArticle struct {
	SN            int64             `gorm:"column:sn" json:"sn,string"`
	CreateTime    string            `json:"create_time"`
	UpdateTime    string            `json:"update_time"`
	Title         string            `json:"title"`          // 文章标题
	Abstract      string            `json:"abstract"`       // 文章简介
	Content       string            `json:"content"`        // 文章内容
	LookCount     int64             `json:"look_count"`     // 浏览量
	CommentCount  int64             `json:"comment_count"`  // 评论量
	DiggCount     int64             `json:"digg_count"`     // 点赞量
	CollectsCount int64             `json:"collects_count"` // 收藏量
	Tags          ctype.ArrayString `json:"tags"`           // 文章标签

	BannerSN  int64  `json:"banner_sn"` // 文章封面 SN
	BannerUrl string `json:"banner_url"`

	UserSN     int64  `json:"user_sn"` // 用户 SN
	Username   string `json:"username"`
	UserAvatar string `json:"user_avatar"`

	IsCollect bool `gorm:"-" json:"is_collect"` // 是否收藏，不入库
	IsDigg    bool `gorm:"-" json:"is_digg"`    // 是否点赞，不入库
}

type ResponseArticleList struct {
	List  []ResponseArticle `json:"list"`
	Count int64             `json:"count"`
}

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
