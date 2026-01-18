package logic

import (
	"github.com/bramble555/blog/dao/mysql/article"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/comment"
	"github.com/bramble555/blog/dao/mysql/digg"
)

func PostArticleDig(uSN, sn int64) (string, error) {
	// 查询文章 sn 是否存在
	ok, err := article.CheckSNExist(sn)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExist
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

// PostArticleCommentDigg 点赞评论或者取消点赞评论，返回点赞状态和错误信息
// true 是现在点赞了, false 是现在取消点赞了
func PostArticleCommentDigg(uSN, sn int64) (bool, error) {
	// 查询 sn 是否存在
	ok, err := comment.CheckSNExist(sn)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, code.ErrorSNNotExist
	}
	isDigg, err := digg.PostCommentDig(uSN, sn)
	if err != nil {
		return false, err
	}
	return isDigg, nil
}
