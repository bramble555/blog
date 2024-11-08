package model

import (
	"context"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model/ctype"
	"github.com/olivere/elastic/v7"
)

type ArticleModel struct {
	ID            uint        `gorm:"primaryKey,autoIncrement" json:"id,string"`
	CreateTime    string      `json:"create_time"`
	UpdateTime    string      `json:"-"`
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
	Tags          ctype.Array `json:"tags"`           // 文章标签

	BannerID  uint   `json:"banner_id,string"` // 文章封面 ID
	BannerUrl string `json:"banner_url"`

	UserID     uint   `json:"user_id,string"` // 用户 ID
	Username   string `json:"username"`
	UserAvatar string `json:"user_avatar"`
}
type ParamArticle struct {
	Title    string      `json:"title" binding:"required"`   // 文章标题
	Abstract string      `json:"abstract"`                   // 文章简介
	Content  string      `json:"content" binding:"required"` // 文章内容
	Category string      `json:"category"`                   // 文章分类
	Source   string      `json:"source"`                     // 文章来源
	Link     string      `json:"link"`                       // 原文链接
	Tags     ctype.Array `json:"tags"`                       // 文章标签
	BannerID uint        `json:"banner_id,string"`           // 文章封面 ID
}
type ResponseArticle struct {
	ID            string      `json:"id,string"` // ES 中 ID
	CreateTime    string      `json:"create_time"`
	UpdateTime    string      `json:"update_time"`
	Title         string      `json:"title"`          // 文章标题
	Abstract      string      `json:"abstract"`       // 文章简介
	LookCount     uint        `json:"look_count"`     // 浏览量
	CommentCount  uint        `json:"comment_count"`  // 评论量
	DiggCount     uint        `json:"digg_count"`     // 点赞量
	CollectsCount uint        `json:"collects_count"` // 收藏量
	Category      string      `json:"category"`       // 文章分类
	Tags          ctype.Array `json:"tags"`           // 文章标签

	BannerID  uint   `json:"banner_id,string"` // 文章封面 ID
	BannerUrl string `json:"banner_url"`

	UserID     uint   `json:"user_id,string"` // 用户 ID
	Username   string `json:"username"`
	UserAvatar string `json:"user_avatar"`
}
type CalendarCount struct {
	Date  string `json:"data"`
	Count int    `json:"count"`
}

// IsExistTitle 判断 title 是否存在
func (a ArticleModel) IsExistTitle(title string) bool {
	boolSearch := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("title", title))
	res, err := global.ES.Search(a.Index()).Query(boolSearch).
		Size(1).
		Do(context.Background())
	if err != nil {
		global.Log.Errorf("IsExistTitle err:%s\n", err.Error())
		return true
	}
	if res.Hits.TotalHits.Value > 0 {
		return true
	}
	return false
}
func (a *ArticleModel) CreateIndex() error {
	exist := a.IndexExists()
	// 创建索引
	if exist {
		return nil
	}
	res, err := global.ES.CreateIndex(a.Index()).BodyString(a.Mapping()).Do(context.Background())
	if !res.Acknowledged {
		global.Log.Printf("创建失败\n")
		return code.ErrorCreateWrong
	}
	return err
}

func (a *ArticleModel) DeleteIndex() error {
	exist := a.IndexExists()
	// 删除索引
	if !exist {
		return nil
	}
	_, err := global.ES.DeleteIndex(a.Index()).Do(context.Background())
	if err != nil {
		global.Log.Errorf("es.DeleteIndex(d.Index()) err:%s\n", err.Error())
		return err
	}
	return nil
}
func (a *ArticleModel) IndexExists() bool {
	index := a.Index()
	exist, err := global.ES.IndexExists(index).Do(context.Background())
	if err != nil {
		global.Log.Errorf("judge IndexExists err:%s\n", err.Error())
		return true
	}
	return exist
}

// Index 返回索引名称
func (ArticleModel) Index() string {
	return "article_index"
}
func (ArticleModel) Mapping() string {
	return `
{
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "create_time": {
         "type": "date",
        "format": "yyyy-MM-dd HH:mm:ss||epoch_millis"
      },
      "update_time": {
         "type": "date",
        "format": "yyyy-MM-dd HH:mm:ss||epoch_millis"
      },
      "title": {
        "type": "text"
      },
      "abstract": {
        "type": "text"
      },
      "content": {
        "type": "text"
      },
      "look_count": {
        "type": "integer"
      },
      "comment_count": {
        "type": "integer"
      },
      "digg_count": {
        "type": "integer"
      },
      "collects_count": {
        "type": "integer"
      },
      "category": {
        "type": "keyword"
      },
      "source": {
        "type": "text"
      },
      "link": {
        "type": "text"
      },
      "tags": {
        "type": "text"
      },
      "banner_id": {
        "type": "integer"
      },
      "banner_url": {
        "type": "text"
      },
      "user_id": {
        "type": "integer"
      },
      "username": {
        "type": "text"
      },
      "user_avatar": {
        "type": "text"
      }
    }
  }
}
`
}
