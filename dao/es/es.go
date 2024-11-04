package es

import (
	"github.com/olivere/elastic/v7"
)

// Init 初始化 Elasticsearch 客户端
func Init() (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
		// 如果需要认证，请设置正确的用户名和密码
		// elastic.SetBasicAuth("username", "password"),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
