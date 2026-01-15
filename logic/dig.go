package logic

import (
	"github.com/bramble555/blog/dao/mysql/article"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/comment"
	"github.com/bramble555/blog/dao/mysql/digg"
)

func PostArticleDig(uSN, sn int64) (string, error) {
	// 查询 sn 是否存在
	ok, err := article.CheckSNExist(sn)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExit
	}
	isDigg, err := digg.PostArticleDig(uSN, sn)
	if err != nil {
		return "", err
	}
	if isDigg {
		return "点赞成功", nil
	}
	return "取消点赞成功", nil
}

func PostArticleCommentDig(uSN, sn int64) (string, error) {
	// 查询 sn 是否存在
	ok, err := comment.CheckSNExist(sn)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExit
	}
	isDigg, err := digg.PostCommentDig(uSN, sn)
	if err != nil {
		return "", err
	}
	if isDigg {
		return "点赞成功", nil
	}
	return "取消点赞成功", nil
}
