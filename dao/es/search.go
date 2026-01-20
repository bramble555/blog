package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/bramble555/blog/global"
)

// esArticle ES 中存储的文章结构
type esArticle struct {
	SN         int64    `json:"sn"`
	Title      string   `json:"title"`
	Tags       []string `json:"tags"`
	Content    string   `json:"content"`
	CreateTime string   `json:"create_time"`
}

// searchResult 搜索结果
type searchResult struct {
	Articles   []articleWithHighlight
	TotalCount int64
}

// articleWithHighlight 带高亮的文章
type articleWithHighlight struct {
	SN               int64  `json:"sn"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	CreateTime       string `json:"create_time"`
	HighlightTitle   string `json:"highlight_title"`   // 高亮后的标题
	HighlightContent string `json:"highlight_content"` // 高亮后的内容
}

// searchArticlesByES 使用 ES 搜索文章
// keyword: 搜索关键词
// page: 页码（从 1 开始）
// size: 每页数量
func searchArticlesByES(keyword string, page, size int) (*searchResult, error) {
	if keyword == "" {
		return &searchResult{Articles: []articleWithHighlight{}, TotalCount: 0}, nil
	}

	// 计算 offset
	from := (page - 1) * size

	// 构建搜索查询
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"title", "content"},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{
					"pre_tags":  []string{`<em style="color:red">`},
					"post_tags": []string{"</em>"},
				},
				"content": map[string]interface{}{
					"pre_tags":  []string{`<em style="color:red">`},
					"post_tags": []string{"</em>"},
				},
			},
		},
		"from": from,
		"size": size,
		"sort": []map[string]interface{}{
			{"_score": map[string]string{"order": "desc"}},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		global.Log.Errorf("SearchArticlesByES encode query failed: %v", err)
		return nil, err
	}

	// 执行搜索
	res, err := global.ES.Search(
		global.ES.Search.WithContext(context.Background()),
		global.ES.Search.WithIndex("articles"),
		global.ES.Search.WithBody(&buf),
		global.ES.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		global.Log.Errorf("SearchArticlesByES search failed: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		global.Log.Errorf("SearchArticlesByES search error: %s", res.String())
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		global.Log.Errorf("SearchArticlesByES decode response failed: %v", err)
		return nil, err
	}

	// 提取命中数据
	hits := result["hits"].(map[string]interface{})
	total := hits["total"].(map[string]interface{})
	totalCount := int64(total["value"].(float64))

	articles := make([]articleWithHighlight, 0)
	hitsList := hits["hits"].([]interface{})

	for _, hit := range hitsList {
		h := hit.(map[string]interface{})
		source := h["_source"].(map[string]interface{})

		article := articleWithHighlight{
			SN:         int64(source["sn"].(float64)),
			Title:      source["title"].(string),
			Content:    source["content"].(string),
			CreateTime: source["create_time"].(string),
		}

		// 提取高亮
		if highlight, ok := h["highlight"].(map[string]interface{}); ok {
			if titleHL, ok := highlight["title"].([]interface{}); ok && len(titleHL) > 0 {
				article.HighlightTitle = titleHL[0].(string)
			} else {
				article.HighlightTitle = article.Title
			}

			if contentHL, ok := highlight["content"].([]interface{}); ok && len(contentHL) > 0 {
				article.HighlightContent = contentHL[0].(string)
			} else {
				article.HighlightContent = article.Content
			}
		} else {
			article.HighlightTitle = article.Title
			article.HighlightContent = article.Content
		}

		articles = append(articles, article)
	}

	return &searchResult{
		Articles:   articles,
		TotalCount: totalCount,
	}, nil
}

// GetArticleSNsByKeyword 根据关键词获取文章 SN 列表
// 用于集成到现有的 MySQL 查询流程中
func GetArticleSNsByKeyword(keyword string) ([]int64, error) {
	// 搜索所有匹配的文章（最多返回 1000 条）
	result, err := searchArticlesByES(keyword, 1, 1000)
	if err != nil {
		return nil, err
	}

	sns := make([]int64, 0, len(result.Articles))
	for _, article := range result.Articles {
		sns = append(sns, article.SN)
	}

	return sns, nil
}

// GetArticleSNsByTag 根据 Tag 搜索文章 SN 列表
func GetArticleSNsByTag(tag string) ([]int64, error) {
	// 构建搜索查询
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"tags": tag,
			},
		},
		"_source": []string{"sn"},
		"size":    1000,
	}

	articles, err := executeQuery(query)
	if err != nil {
		return nil, err
	}

	sns := make([]int64, 0, len(articles))
	for _, article := range articles {
		sns = append(sns, article.SN)
	}
	return sns, nil
}

func executeQuery(query map[string]interface{}) ([]esArticle, error) {
	// Helper function to execute query and return simplified esArticle list (with SN)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := global.ES.Search(
		global.ES.Search.WithContext(context.Background()),
		global.ES.Search.WithIndex("articles"),
		global.ES.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})
	hitsList := hits["hits"].([]interface{})

	articles := make([]esArticle, 0)
	for _, hit := range hitsList {
		h := hit.(map[string]interface{})
		source := h["_source"].(map[string]interface{})
		sn := int64(source["sn"].(float64))
		articles = append(articles, esArticle{SN: sn})
	}
	return articles, nil
}
