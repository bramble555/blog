package es

import (
	"context"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func UploadArticles(article model.ArticleModel) (string, error) {
	_, err := global.ES.Index().Index(article.Index()).BodyJson(article).Do(context.Background())
	if err != nil {
		global.Log.Errorf("es UploadArticles err:%s\n", err.Error())
		return "", nil
	}
	return code.CreateSucceed, nil
}
