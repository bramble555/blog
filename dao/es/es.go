package es

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/olivere/elastic/v7"
)

// Init 初始化 Elasticsearch 客户端
func Init() (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%d", global.Config.ES.Host, global.Config.ES.Port)),
		elastic.SetSniff(global.Config.ES.Sniff),
		// 如果需要认证，请设置正确的用户名和密码
		// elastic.SetBasicAuth("username", "password"),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
