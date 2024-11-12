package logic

import (
	"github.com/bramble555/blog/dao/mysql/article"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/comment"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func PostArticleComments(uID uint, pc *model.ParamPostComment) (string, error) {
	// 判断文章是否存在
	ok, err := article.IDExist(pc.ArticleID)
	if err != nil {
		return "", err
	}
	if !ok {
		global.Log.Errorf("文章:%d不存在", pc.ArticleID)
		return "", code.ErrorIDNotExit
	}
	if pc.ParentCommentID == -1 {
		return comment.PostArticleComments(uID, pc)
	}
	// 查看父评论是否存在
	ok, err = comment.CheckIDExist(uint(pc.ParentCommentID))
	if err != nil {
		return "", err
	}
	if !ok {
		global.Log.Errorf("评论:%d不存在", pc.ArticleID)
		return "", code.ErrorIDNotExit
	}
	return comment.PostArticleComments(uID, pc)
}
func GetArticleComments(pcl *model.ParamCommentList) ([]model.ResponseCommentList, error) {
	// 判断文章是否存在
	ok, err := article.IDExist(pcl.ArticleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		global.Log.Errorf("文章:%d不存在", pcl.ArticleID)
		return nil, code.ErrorIDNotExit
	}
	return comment.GetArticleComments(pcl)
}
func DeleteArticleComments(uID uint, pi *model.ParamID) (string, error) {
	// 检查 id 是否存在
	ok, err := comment.CheckIDExist(pi.ID)
	if err != nil {
		return "", err
	}
	if !ok {
		global.Log.Errorf("id:%d不存在", pi.ID)
		return "", code.ErrorIDExit
	}
	articleID, err := comment.GetArticleIDByID(pi.ID)
	if err != nil {
		return "", err
	}
	return comment.DeleteArticleComments(uID, pi, articleID)
}
