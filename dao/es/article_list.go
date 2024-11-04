package es

import (
	"context"
	"encoding/json"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/olivere/elastic/v7"
)

func GetArticlesList(pl *model.ParamList) (*[]model.ResponseArticle, error) {
	a := model.ArticleModel{}
	// 设置分页参数
	page := (pl.Page - 1) * pl.Size
	// 构建搜索请求
	resp, err := global.ES.Search().
		Index(a.Index()).Query(elastic.NewBoolQuery()).
		FetchSourceContext(elastic.NewFetchSourceContext(true)).
		From(page).
		Size(pl.Size).
		Do(context.Background())
	if err != nil {
		global.Log.Errorf("search document failed, err:%v\n", err)
		return nil, err
	}

	// 解析结果
	var articles []model.ResponseArticle
	for _, hit := range resp.Hits.Hits {
		var article model.ResponseArticle
		if err := json.Unmarshal(hit.Source, &article); err != nil {
			global.Log.Errorf("failed to unmarshal document: %v\n", err)
			continue
		}
		articles = append(articles, article)
	}
	return &articles, nil
}
