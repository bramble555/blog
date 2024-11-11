package logic

import (
	"github.com/bramble555/blog/dao/mysql/article"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/comment"
	"github.com/bramble555/blog/dao/redis"
)

func PostArticleDig(id uint) (string, error) {
	// 查询 id 是否存在
	ok, err := article.IDExist(id)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorIDNotExit
	}
	return redis.PostArticleDig(id)
}
func PostArticleCommentDig(id uint) (string, error) {
	// 查询 id 是否存在
	ok, err := comment.IDExist(id)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorIDNotExit
	}
	return redis.PostArticleCommentDig(id)
}
