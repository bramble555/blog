package logic

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/article"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/comment"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func PostArticleComments(uSN int64, pc *model.ParamPostComment) (string, error) {
	// 判断文章是否存在
	ok, err := article.CheckSNExist(pc.ArticleSN)
	if err != nil {
		return "", err
	}
	if !ok {
		global.Log.Errorf("文章:%d不存在", pc.ArticleSN)
		return "", code.ErrorSNNotExit
	}
	// 判断父评论是否存在
	if pc.ParentCommentSN != -1 {
		ok, err := comment.CheckSNExist(pc.ParentCommentSN)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", code.ErrorSNNotExit
		}
		// 获取父评论的文章 SN
		articleSN, err := comment.GetArticleSNBySN(pc.ParentCommentSN)
		if err != nil {
			return "", err
		}
		// 校验文章 SN 是否一致
		if articleSN != pc.ArticleSN {
			return "", errors.New("评论的文章 SN 与父评论所属文章 SN 不一致")
		}
	}
	return comment.PostArticleComments(uSN, pc)
}
func GetArticleComments(pcl *model.ParamCommentList, uSN int64) (*model.ResponseCommentListWrapper, error) {
	// 如果未指定文章 SN，则返回全站评论（后台用途）
	if pcl.ArticleSN == 0 {
		return comment.GetAllComments(pcl, uSN)
	}
	// 否则校验文章存在并返回该文章下的评论
	ok, err := article.CheckSNExist(pcl.ArticleSN)
	if err != nil {
		return nil, err
	}
	if !ok {
		global.Log.Errorf("文章:%d不存在", pcl.ArticleSN)
		return nil, code.ErrorSNNotExit
	}
	return comment.GetArticleComments(pcl, uSN)
}
func DeleteArticleComments(uSN int64, role int64, psn *model.ParamSN, articleSN int64) (string, error) {
	// 如果 articleSN 为 0，说明是从控制器传来的，需要查询
	if articleSN == 0 {
		var err error
		articleSN, err = comment.GetArticleSNBySN(psn.SN)
		if err != nil {
			return "", err
		}
	}
	return comment.DeleteArticleComments(uSN, role, psn, articleSN)
}
