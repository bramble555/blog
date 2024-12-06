package es

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/olivere/elastic/v7"
)

type ESArticleQueryService struct{}

func (e *ESArticleQueryService) GetArticlesListByParam(paq *model.ParamArticleQuery) (*[]model.ResponseArticle, error) {
	a := model.ArticleModel{}
	// 设置分页参数
	page := (paq.Page - 1) * paq.Size

	// 构建搜索请求的查询部分
	boolQuery := elastic.NewBoolQuery()

	// 如果传递了 title 参数，构造全文搜索条件
	if paq.Title != "" {
		boolQuery.Must(elastic.NewMatchQuery("title", paq.Title))
	}

	// 如果传递了 tags 参数
	if paq.Tags != "" {
		boolQuery.Must(elastic.NewMatchQuery("tags", paq.Tags))
	}
	if paq.Content != "" {
		boolQuery.Must(elastic.NewMatchQuery("content", paq.Content))
	}

	// 构建搜索请求
	searchService := global.ES.Search().
		Index(a.Index()).
		Query(boolQuery).
		FetchSourceContext(elastic.NewFetchSourceContext(true)).
		From(page).
		Size(paq.Size)

	// 如果 title 参数存在，则启用高亮功能
	if paq.Title != "" {
		searchService.Highlight(elastic.NewHighlight().Field("title"))
	}

	// 执行查询
	resp, err := searchService.Do(context.Background())
	if err != nil {
		global.Log.Errorf("search document failed, err:%v\n", err)
		return nil, err
	}

	// 解析结果
	var articles []model.ResponseArticle
	for _, hit := range resp.Hits.Hits {
		var article model.ResponseArticle
		// 如果是单一的 string，也就是只有一个 tag
		if err := json.Unmarshal(hit.Source, &article); err != nil {
			// 出错了，说明这个文章对应多个 tag
			// 因为 tags 存储数据库的时候，一对多的时候 以 `\n` 分割
			// ES 与 MySQL 是同步的，所以 tags 在 ES 里面也是 有 `\n`的
			// 尝试分割 `\n`
			var rawMap map[string]interface{}
			if err := json.Unmarshal(hit.Source, &rawMap); err == nil {
				if tags, ok := rawMap["tags"]; ok {
					switch v := tags.(type) {
					case string:
						article.Tags = strings.Split(v, "\n")
					default:
						global.Log.Errorf("unknown type for tags: %T", v)
					}
				}
			}
		}
		// 如果有高亮信息，填充到对应字段中
		if len(hit.Highlight["title"]) > 0 {
			article.Title = hit.Highlight["title"][0]
		}

		articles = append(articles, article)
	}

	return &articles, nil
}

func GetArticlesDetail(id string) (*model.ArticleModel, error) {
	a := model.ArticleModel{}
	res, err := global.ES.Get().Index(model.ArticleModel{}.Index()).Id(id).Do(context.Background())
	if err != nil {
		global.Log.Errorf("global.ES.Get() err:%s\n", err.Error())
		return nil, err
	}
	err = json.Unmarshal(res.Source, &a)
	if err != nil {
		global.Log.Errorf("json.Unmarshal err:%s\n", err.Error())
		return nil, err
	}
	return &a, nil
}
