package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

// SyncArticlesToES 每天晚上0点同步文章到 ES（增量同步）
// 新增：从 MySQL 读取所有文章并同步到 ES
// 删除：删除 ES 中存在但 MySQL 中不存在的文章
func SyncArticlesToES() error {
	global.Log.Info("Starting incremental sync articles to ES...")

	if global.ES == nil {
		global.Log.Error("ES client is nil, sync aborted")
		return fmt.Errorf("ES client is nil")
	}

	// 1. 已在 InitES 中完成索引检查与创建

	// 2. 从 MySQL 获取所有文章的 SN
	var mysqlSNs []int64
	if err := global.DB.Table("article_models").
		Pluck("sn", &mysqlSNs).Error; err != nil {
		global.Log.Errorf("SyncArticlesToES query MySQL SNs failed: %v", err)
		return err
	}

	mysqlSNMap := make(map[int64]bool)
	for _, sn := range mysqlSNs {
		mysqlSNMap[sn] = true
	}

	global.Log.Infof("Found %d articles in MySQL", len(mysqlSNs))

	// 3. 从 ES 获取所有文章的 SN
	esSNs, err := getAllArticleSNsFromES()
	if err != nil {
		global.Log.Errorf("SyncArticlesToES get ES SNs failed: %v", err)
		return err
	}

	global.Log.Infof("Found %d articles in ES", len(esSNs))

	// 4. 找出需要删除的文章（ES 中存在但 MySQL 中不存在）
	var toDelete []int64
	for _, sn := range esSNs {
		if !mysqlSNMap[sn] {
			toDelete = append(toDelete, sn)
		}
	}

	if len(toDelete) > 0 {
		global.Log.Infof("Deleting %d orphaned articles from ES", len(toDelete))
		if err := deleteArticlesFromES(toDelete); err != nil {
			global.Log.Errorf("SyncArticlesToES delete orphaned articles failed: %v", err)
			return err
		}
	}

	// 5. 同步所有 MySQL 中的文章到 ES（Bulk API）
	// 这样可以确保新增和更新的文章都会被同步
	var articles []model.ArticleModel
	if err := global.DB.Table("article_models").
		Select("sn, title, content, create_time, tags").
		Find(&articles).Error; err != nil {
		global.Log.Errorf("SyncArticlesToES query articles failed: %v", err)
		return err
	}

	if len(articles) == 0 {
		global.Log.Info("No articles to sync")
		return nil
	}

	// 6. 使用 Bulk API 批量写入/更新 ES
	batchSize := 100
	for i := 0; i < len(articles); i += batchSize {
		end := i + batchSize
		if end > len(articles) {
			end = len(articles)
		}

		batch := articles[i:end]
		if err := bulkIndexArticles(batch); err != nil {
			global.Log.Errorf("SyncArticlesToES bulk index failed at batch %d-%d: %v", i, end, err)
			return err
		}

		global.Log.Infof("Synced articles %d-%d", i+1, end)
	}

	global.Log.Infof("Incremental sync completed successfully")
	global.Log.Infof("Total: %d articles in MySQL, %d deleted from ES", len(articles), len(toDelete))
	return nil
}

// bulkIndexArticles 批量索引文章
func bulkIndexArticles(articles []model.ArticleModel) error {
	var buf bytes.Buffer

	for _, article := range articles {
		// Bulk API 格式：{ "index": { "_index": "articles", "_id": "sn" } }
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "articles",
				"_id":    fmt.Sprintf("%d", article.SN),
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}

		// 文档数据
		doc := esArticle{
			SN:         article.SN,
			Title:      article.Title,
			Content:    article.Content,
			Tags:       article.Tags,
			CreateTime: article.CreateTime.Format("2006-01-02 15:04:05"),
		}
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return err
		}
	}

	// 执行 Bulk 请求
	res, err := global.ES.Bulk(
		bytes.NewReader(buf.Bytes()),
		global.ES.Bulk.WithContext(context.Background()),
		global.ES.Bulk.WithIndex("articles"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var raw map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
			return fmt.Errorf("bulk error: %s", res.String())
		}

		// 检查是否有错误
		if raw["errors"].(bool) {
			items := raw["items"].([]interface{})
			var errMsgs []string
			for _, item := range items {
				itemMap := item.(map[string]interface{})
				if indexResult, ok := itemMap["index"].(map[string]interface{}); ok {
					if indexResult["error"] != nil {
						errMsgs = append(errMsgs, fmt.Sprintf("%v", indexResult["error"]))
					}
				}
			}
			return fmt.Errorf("bulk errors: %s", strings.Join(errMsgs, "; "))
		}
	}

	return nil
}

// getAllArticleSNsFromES 获取 ES 中所有文章的 SN
func getAllArticleSNsFromES() ([]int64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"_source": []string{"sn"},
		"size":    10000, // 最多获取 10000 条
	}

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

	sns := make([]int64, 0, len(hitsList))
	for _, hit := range hitsList {
		h := hit.(map[string]interface{})
		source := h["_source"].(map[string]interface{})
		sn := int64(source["sn"].(float64))
		sns = append(sns, sn)
	}

	return sns, nil
}

// deleteArticlesFromES 从 ES 批量删除文章
func deleteArticlesFromES(sns []int64) error {
	if len(sns) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, sn := range sns {
		// Bulk delete 格式
		meta := map[string]interface{}{
			"delete": map[string]interface{}{
				"_index": "articles",
				"_id":    fmt.Sprintf("%d", sn),
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}
	}

	res, err := global.ES.Bulk(
		bytes.NewReader(buf.Bytes()),
		global.ES.Bulk.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk delete error: %s", res.String())
	}

	global.Log.Infof("Deleted %d articles from ES", len(sns))
	return nil
}

// UploadIndexSingleArticle 上传文章的时候,索引单篇文章到 ES（实时同步）
func UploadIndexSingleArticle(articleSN int64) error {
	if global.ES == nil {
		global.Log.Warn("ES client is nil, skip indexing single article")
		return nil
	}

	// 从 MySQL 查询文章信息
	var article model.ArticleModel
	if err := global.DB.Table("article_models").
		Select("sn, title, content, create_time, tags").
		Where("sn = ?", articleSN).
		First(&article).Error; err != nil {
		global.Log.Errorf("IndexSingleArticle query article failed: %v", err)
		return err
	}

	// 构建文档
	doc := esArticle{
		SN:         article.SN,
		Title:      article.Title,
		Content:    article.Content,
		Tags:       article.Tags,
		CreateTime: article.CreateTime.Format("2006-01-02 15:04:05"),
	}

	// 转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		global.Log.Errorf("IndexSingleArticle encode doc failed: %v", err)
		return err
	}

	// 索引到 ES
	res, err := global.ES.Index(
		"articles",
		&buf,
		global.ES.Index.WithDocumentID(fmt.Sprintf("%d", articleSN)),
		global.ES.Index.WithContext(context.Background()),
	)
	if err != nil {
		global.Log.Errorf("IndexSingleArticle index failed: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		global.Log.Errorf("IndexSingleArticle index error: %s", res.String())
		return fmt.Errorf("index error: %s", res.String())
	}

	global.Log.Infof("Article %d indexed to ES successfully", articleSN)
	return nil
}
