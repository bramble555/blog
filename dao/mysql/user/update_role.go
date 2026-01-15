package user

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

func UpdateUserRole(puur *model.ParamUpdateUserRole) (string, error) {
	err := global.DB.Table("user_models").Where("sn = ?", puur.UserSN).Update("role", puur.Role).Error
	if err != nil {
		global.Log.Errorf("user UpdateUserRole err:%s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("修改用户%d权限成功", puur.UserSN), nil
}

// UpdateUserPwd 负责更新用户密码
func UpdateUserPwd(puup *model.ParamUpdateUserPwd, sn int64) (string, error) {
	// 先把先密码加密
	pwd, err := pkg.HashPassword(puup.Pwd)
	if err != nil {
		global.Log.Errorf("HashPassword err:%s\n", err.Error())
		return "", err
	}
	err = global.DB.Table("user_models").Where("sn = ?", sn).
		Update("password", pwd).Error
	if err != nil {
		global.Log.Errorf("user UpdateUserPwd err: %s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("修改用户 %d 密码成功", sn), nil
}

func DeleteUser(sn int64) (string, error) {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户关联的文章
	// var articleSNs []int64
	// if err := tx.Table("article_models").Where("user_sn = ?", sn).Pluck("sn", &articleSNs).Error; err != nil {
	// 	tx.Rollback()
	// 	return "", err
	// }
	// if len(articleSNs) > 0 {
	// 	if err := tx.Table("comment_models").Where("article_sn IN ?", articleSNs).Delete(&model.CommentModel{}).Error; err != nil {
	// 		tx.Rollback()
	// 		return "", err
	// 	}
	// 	if err := tx.Table("article_tag_models").Where("article_sn IN ?", articleSNs).Delete(&model.ArticleTagModel{}).Error; err != nil {
	// 		tx.Rollback()
	// 		return "", err
	// 	}
	// 	if err := tx.Table("article_models").Where("sn IN ?", articleSNs).Delete(&model.ArticleModel{}).Error; err != nil {
	// 		tx.Rollback()
	// 		return "", err
	// 	}
	// }

	// 删除用户关联的评论
	if err := tx.Table("comment_models").Where("user_sn = ?", sn).Delete(&model.CommentModel{}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// 删除用户
	if err := tx.Table("user_models").Where("sn = ?", sn).Delete(&model.UserModel{}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("删除用户 %d 及其文章和评论成功", sn), nil
}
