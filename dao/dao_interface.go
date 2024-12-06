package dao

import "github.com/bramble555/blog/model"

type ArticleQueryService interface {
	GetArticlesListByParam(paq *model.ParamArticleQuery) (*[]model.ResponseArticle, error)
}

