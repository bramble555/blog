package digg

import (
	"errors"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"gorm.io/gorm"
)

// PostArticleDig 文章点赞/取消点赞
// 返回: true = 点赞成功, false = 取消点赞成功, error
func PostArticleDig(uSN, articleSN int64) (bool, error) {
	var userDigg model.UserDiggModel
	err := global.DB.Where("user_sn = ? AND article_sn = ?", uSN, articleSN).First(&userDigg).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 未点赞 -> 点赞
		return true, global.DB.Transaction(func(tx *gorm.DB) error {
			// 1. 创建点赞记录
			if err := tx.Create(&model.UserDiggModel{UserSN: uSN, ArticleSN: articleSN}).Error; err != nil {
				return err
			}
			// 2. 文章点赞数 +1
			if err := tx.Model(&model.ArticleModel{}).Where("sn = ?", articleSN).Update("digg_count", gorm.Expr("digg_count + ?", 1)).Error; err != nil {
				return err
			}
			return nil
		})
	} else if err == nil {
		// 已点赞 -> 取消点赞
		return false, global.DB.Transaction(func(tx *gorm.DB) error {
			// 1. 删除点赞记录
			if err := tx.Where("user_sn = ? AND article_sn = ?", uSN, articleSN).Delete(&model.UserDiggModel{}).Error; err != nil {
				return err
			}
			// 2. 文章点赞数 -1
			if err := tx.Model(&model.ArticleModel{}).Where("sn = ?", articleSN).Update("digg_count", gorm.Expr("digg_count - ?", 1)).Error; err != nil {
				return err
			}
			return nil
		})
	} else {
		return false, err
	}
}

// PostCommentDig 评论点赞/取消点赞
func PostCommentDig(uSN, commentSN int64) (bool, error) {
	var userDigg model.UserCommentDiggModel
	err := global.DB.Where("user_sn = ? AND comment_sn = ?", uSN, commentSN).First(&userDigg).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 未点赞 -> 点赞
		return true, global.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&model.UserCommentDiggModel{UserSN: uSN, CommentSN: commentSN}).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.CommentModel{}).Where("sn = ?", commentSN).Update("digg_count", gorm.Expr("digg_count + ?", 1)).Error; err != nil {
				return err
			}
			return nil
		})
	} else if err == nil {
		// 已点赞 -> 取消点赞
		return false, global.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("user_sn = ? AND comment_sn = ?", uSN, commentSN).Delete(&model.UserCommentDiggModel{}).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.CommentModel{}).Where("sn = ?", commentSN).Update("digg_count", gorm.Expr("digg_count - ?", 1)).Error; err != nil {
				return err
			}
			return nil
		})
	} else {
		return false, err
	}
}

func IsUserDigg(uSN, articleSN int64) (bool, error) {
	var count int64
	err := global.DB.Model(&model.UserDiggModel{}).Where("user_sn = ? AND article_sn = ?", uSN, articleSN).Count(&count).Error
	return count > 0, err
}

func IsUserCommentDigg(uSN, commentSN int64) (bool, error) {
	var count int64
	err := global.DB.Model(&model.UserCommentDiggModel{}).Where("user_sn = ? AND comment_sn = ?", uSN, commentSN).Count(&count).Error
	return count > 0, err
}
