package logic

import (
	"github.com/bramble555/blog/dao/mysql/article"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/comment"
	"github.com/bramble555/blog/dao/redis"
)

func PostArticleDig(sn int64) (string, error) {
	// 查询 sn 是否存在
	ok, err := article.CheckSNExist(sn)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExit
	}
	return redis.PostArticleDig(sn)
}
func PostArticleCommentDig(sn int64) (string, error) {
	// 查询 sn 是否存在
	ok, err := comment.CheckSNExist(sn)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExit
	}
	return redis.PostArticleCommentDig(sn)
}
