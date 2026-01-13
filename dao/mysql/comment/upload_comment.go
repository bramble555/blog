package comment

import (
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"gorm.io/gorm"
)

func CheckSNExist(sn int64) (bool, error) {
	var count int64
	err := global.DB.Table("comment_models").Where("sn = ?", sn).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error SNExist: %v\n", err)
		return false, code.ErrorSNNotExit
	}
	return count > 0, nil
}
func PostArticleComments(uSN int64, pc *model.ParamPostComment) (string, error) {
	tx := global.DB.Begin()
	// 把此次评论添加到 comment_models 中
	err := tx.Table("comment_models").Create(&model.CommentModel{
		Content:         pc.Content,
		ParentCommentSN: pc.ParentCommentSN,
		ArticleSN:       pc.ArticleSN,
		UserSN:          uSN,
	}).Error
	if err != nil {
		global.Log.Errorf("comment_models  Create err:%s\n", err.Error())
		tx.Rollback()
		return "", err
	}
	err = tx.Table("article_models").Where("sn = ?", pc.ArticleSN).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).
		Error
	if err != nil {
		global.Log.Errorf("scan comment_count err:%s\n", err.Error())
		tx.Rollback()
		return "", err
	}
	// 更新父级评论的 CommentCount
	if err = updateParentCommentCount(tx, pc.ParentCommentSN); err != nil {
		global.Log.Errorf("updateParentCommentCount err:%s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("Commit err:%s\n", err.Error())
		tx.Rollback()
		return "", err
	}
	return "评论成功", nil
}
func GetArticleSNBySN(sn int64) (int64, error) {
	var articleSN int64
	err := global.DB.Table("comment_models").Where("sn = ?", sn).
		Select("article_sn").Scan(&articleSN).Error
	if err != nil {
		global.Log.Errorf("Error SNExist: %v\n", err)
		return 0, err
	}
	return articleSN, nil
}

// updateParentCommentCount 递归更新父级评论的 CommentCount
// 从下到上
func updateParentCommentCount(tx *gorm.DB, parentSN int64) error {
	// 如果没有父级评论则直接返回
	if parentSN == -1 {
		return nil
	}

	// 更新当前父级评论的 CommentCount +1
	if err := tx.Table("comment_models").Where("sn = ?", parentSN).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).
		Error; err != nil {
		return err
	}

	// 查询当前父级评论的 ParentCommentSN
	var parentComment model.CommentModel
	if err := tx.Table("comment_models").Where("sn = ?", parentSN).
		First(&parentComment).
		Error; err != nil {
		return err
	}

	// 递归更新上一级父评论
	return updateParentCommentCount(tx, parentComment.ParentCommentSN)
}
