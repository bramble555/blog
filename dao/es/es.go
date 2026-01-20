package es

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/bramble555/blog/global"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

// InitES 初始化 Elasticsearch 客户端
func InitES() (*elasticsearch.Client, error) {
	if !global.Config.Elasticsearch.Enable {
		global.Log.Info("Elasticsearch is disabled by config")
		return nil, nil
	}
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%d", global.Config.Elasticsearch.Host, global.Config.Elasticsearch.Port),
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略证书校验
			},
		},
	}

	// 如果配置了用户名和密码，则添加认证
	if global.Config.Elasticsearch.Username != "" && global.Config.Elasticsearch.Password != "" {
		cfg.Username = global.Config.Elasticsearch.Username
		cfg.Password = global.Config.Elasticsearch.Password
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		global.Log.Errorf("InitES failed: %v", err)
		return nil, err
	}

	// 测试连接
	res, err := client.Info()
	if err != nil {
		global.Log.Errorf("ES connection test failed: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		global.Log.Errorf("ES connection test error: %s", res.String())
		return nil, fmt.Errorf("ES connection error: %s", res.String())
	}

	global.Log.Info("Elasticsearch connected successfully")

	// 将 client 赋值给 global.ES，这样后续的辅助函数可以使用 global.ES
	global.ES = client

	// 检查并创建 articles 索引
	exists, err := indexExists("articles")
	if err != nil {
		global.Log.Errorf("InitES check articles index exists failed: %v", err)
		return client, nil // 即使检查失败，也返回 client，不干扰主线程
	}

	if !exists {
		global.Log.Info("Articles index does not exist, creating in InitES...")
		if err := createArticlesIndex(); err != nil {
			global.Log.Errorf("InitES create articles index failed: %v", err)
		}
	}

	return client, nil
}

// indexExists 检查索引是否存在
func indexExists(indexName string) (bool, error) {
	res, err := global.ES.Indices.Exists([]string{indexName})
	if err != nil {
		global.Log.Errorf("IndexExists check failed: %v", err)
		return false, err
	}
	defer res.Body.Close()

	// 200 表示存在，404 表示不存在
	return res.StatusCode == 200, nil
}

// deleteIndex 删除索引
func deleteIndex(indexName string) error {
	res, err := global.ES.Indices.Delete([]string{indexName})
	if err != nil {
		global.Log.Errorf("deleteIndex failed: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		global.Log.Errorf("deleteIndex error: %s", res.String())
		return fmt.Errorf("delete index error: %s", res.String())
	}

	global.Log.Infof("Index %s deleted successfully", indexName)
	return nil
}

// createArticlesIndex 创建 articles 索引及 Mapping
// 优先使用 ik_max_word 分词器，如果不可用则使用 standard 分词器
func createArticlesIndex() error {
	indexName := "articles"

	// 定义 Mapping，优先使用 ik_max_word，如果不可用则使用 standard
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"default": {
						"type": "standard"
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"sn": {
					"type": "keyword"
				},
				"title": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				},
				"content": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				},
				"tags": {
					"type": "keyword"
				},
				"create_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				}
			}
		}
	}`

	res, err := global.ES.Indices.Create(
		indexName,
		global.ES.Indices.Create.WithBody(strings.NewReader(mapping)),
		global.ES.Indices.Create.WithContext(context.Background()),
	)
	if err != nil {
		global.Log.Errorf("CreateArticlesIndex failed: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		errMsg := res.String()
		// 如果是因为 ik_max_word 分词器不存在，则使用 standard 分词器
		if strings.Contains(errMsg, "ik_max_word") || strings.Contains(errMsg, "unknown analyzer") {
			global.Log.Warnf("ik_max_word analyzer not available, using standard analyzer")
			return createArticlesIndexWithStandard()
		}
		global.Log.Errorf("createArticlesIndex error: %s", errMsg)
		return fmt.Errorf("create articles index error: %s", errMsg)
	}

	global.Log.Info("Articles index created successfully with ik_max_word analyzer")
	return nil
}

// createArticlesIndexWithStandard 使用 standard 分词器创建索引
func createArticlesIndexWithStandard() error {
	indexName := "articles"

	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"sn": {
					"type": "keyword"
				},
				"title": {
					"type": "text",
					"analyzer": "standard"
				},
				"content": {
					"type": "text",
					"analyzer": "standard"
				},
				"tags": {
					"type": "keyword"
				},
				"create_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				}
			}
		}
	}`

	res, err := global.ES.Indices.Create(
		indexName,
		global.ES.Indices.Create.WithBody(strings.NewReader(mapping)),
		global.ES.Indices.Create.WithContext(context.Background()),
	)
	if err != nil {
		global.Log.Errorf("createArticlesIndexWithStandard failed: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		global.Log.Errorf("createArticlesIndexWithStandard error: %s", res.String())
		return fmt.Errorf("create articles index with standard error: %s", res.String())
	}

	global.Log.Info("Articles index created successfully with standard analyzer")
	return nil
}
