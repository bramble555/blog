package comment

import (
	"fmt"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"gorm.io/gorm"
)

func IDExist(id uint) (bool, error) {
	var count int64
	err := global.DB.Table("comment_models").Where("id = ?", id).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error IDExist: %v\n", err)
		return false, code.ErrorIDNotExit
	}
	return count > 0, nil
}
func PostArticleComments(uID uint, pc *model.ParamPostComment) (string, error) {
	tx := global.DB.Begin()
	// 把此次评论添加到 comment_models 中
	err := tx.Table("comment_models").Create(&model.CommentModel{
		Content:         pc.Content,
		ParentCommentID: pc.ParentCommentID,
		ArticleID:       pc.ArticleID,
		UserID:          uID,
	}).Error
	if err != nil {
		global.Log.Errorf("comment_models  Create err:%s\n", err.Error())
		tx.Rollback()
		return "", err
	}
	err = tx.Table("article_models").Where("id = ?", pc.ArticleID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).
		Error
	if err != nil {
		global.Log.Errorf("scan comment_count err:%s\n", err.Error())
		tx.Rollback()
		return "", err
	}
	// 更新父级评论的 CommentCount
	if err = updateParentCommentCount(tx, pc.ParentCommentID); err != nil {
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
func GetArticleComments(pcl *model.ParamCommentList) ([]model.ResponseCommentList, error) {
	// 获取根评论
	rootComments, err := GetRootComments(pcl.ArticleID)
	if err != nil {
		return nil, err
	}
	return rootComments, nil
}

func GetRootComments(articleID uint) ([]model.ResponseCommentList, error) {
	// 查找根评论
	var rootComments []model.CommentModel
	err := global.DB.Table("comment_models").
		Where("article_id = ? AND parent_comment_id = -1", articleID).
		Order("create_time ASC"). // 按创建时间排序
		Find(&rootComments).Error
	if err != nil {
		global.Log.Errorf("comment_models Find err: %s\n", err.Error())
		return nil, err
	}

	// 构建响应数据
	responseCommentsList := make([]model.ResponseCommentList, 0, len(rootComments))
	for _, comment := range rootComments {
		// 查找用户详情，直接通过评论中的 UserID 获取
		userDetail, err := user.GetUserDetailByID(comment.UserID)
		if err != nil {
			return nil, err
		}

		// 递归获取子评论
		subComments, err := getSubComments(comment.ID, articleID)
		if err != nil {
			return nil, err
		}

		// 构建响应评论
		responseComment := model.ResponseCommentList{
			MODEL:           comment.MODEL,
			Content:         comment.Content,
			ParentCommentID: comment.ParentCommentID,
			ArticleID:       comment.ArticleID,
			DiggCount:       comment.DiggCount,
			CommentCount:    comment.CommentCount,
			SubComments:     subComments, // 添加子评论
			UserDetail:      userDetail,
		}

		// 将根评论添加到响应列表
		responseCommentsList = append(responseCommentsList, responseComment)
	}

	return responseCommentsList, nil
}

// getSubComments 递归得到子评论
// 从上到下
func getSubComments(parentCommentID uint, articleID uint) ([]model.ResponseCommentList, error) {
	// 查找子评论
	var subComments []model.CommentModel
	err := global.DB.Table("comment_models").
		Where("article_id = ? AND parent_comment_id = ?", articleID, parentCommentID).
		Order("create_time ASC"). // 按创建时间排序
		Find(&subComments).Error
	if err != nil {
		global.Log.Errorf("comment_models Find err: %s\n", err.Error())
		return nil, err
	}

	// 如果没有子评论，直接返回空
	if len(subComments) == 0 {
		return nil, nil
	}

	// 构建响应数据
	responseSubComments := make([]model.ResponseCommentList, 0, len(subComments))
	for _, comment := range subComments {
		// 查找用户详情，直接通过评论中的 UserID 获取

		userDetail, err := user.GetUserDetailByID(comment.UserID)
		if err != nil {
			// 如果获取不到用户详情，跳过这个评论
			global.Log.Errorf("未找到用户详情, 用户ID: %d", comment.UserID)
			continue
		}

		// 递归获取子评论
		var subSubComments []model.ResponseCommentList
		if comment.CommentCount > 0 { // 只有在有子评论时才递归
			subSubComments, err = getSubComments(comment.ID, articleID)
			if err != nil {
				return nil, fmt.Errorf("获取子评论失败: %w", err)
			}
		}

		// 构建响应评论
		responseComment := model.ResponseCommentList{
			MODEL:           comment.MODEL,
			Content:         comment.Content,
			ParentCommentID: comment.ParentCommentID,
			ArticleID:       comment.ArticleID,
			DiggCount:       comment.DiggCount,
			CommentCount:    comment.CommentCount,
			SubComments:     subSubComments,
			UserDetail:      userDetail,
		}

		// 将子评论添加到响应列表
		responseSubComments = append(responseSubComments, responseComment)
	}

	return responseSubComments, nil
}

// updateParentCommentCount 递归更新父级评论的 CommentCount
// 从下到上
func updateParentCommentCount(tx *gorm.DB, parentID int) error {
	if parentID == -1 {
		return nil // 如果没有父级评论则直接返回
	}

	// 更新当前父级评论的 CommentCount +1
	if err := tx.Table("comment_models").Where("id = ?", parentID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).
		Error; err != nil {
		return err
	}

	// 查询当前父级评论的 ParentCommentID
	var parentComment model.CommentModel
	if err := tx.Table("comment_models").Where("id = ?", parentID).
		First(&parentComment).
		Error; err != nil {
		return err
	}

	// 递归更新上一级父评论
	return updateParentCommentCount(tx, parentComment.ParentCommentID)
}
