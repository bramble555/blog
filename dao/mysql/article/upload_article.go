package article

import (
	"time"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"gorm.io/gorm"
)

func TitleIsExist(title string) (bool, error) {
	var count int64
	err := global.DB.Table("article_models").Where("title = ?", title).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error checking if title exists: %v\n", err)
		return false, code.ErrorTitleExit
	}
	return count > 0, nil
}
func UploadArticles(am *model.ArticleModel) (string, error) {
	// 开始事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Log.Errorf("article  start transaction err:%s\n", tx.Error.Error())
		return "", tx.Error
	}
	// 组装数据,其中 SN 必须现在就要创建了
	now := time.Now()
	am.SN = global.Snowflake.GetID()
	tagsStrList := []string(am.Tags)

	err := tx.Create(&model.ArticleModel{
		MODEL: model.MODEL{
			SN:         am.SN,
			CreateTime: now,
			UpdateTime: now,
		},
		Title:         am.Title,
		Abstract:      am.Abstract,
		Content:       am.Content,
		LookCount:     am.LookCount,
		CommentCount:  am.CommentCount,
		DiggCount:     am.DiggCount,
		CollectsCount: am.CollectsCount,
		Tags:          am.Tags,
		BannerSN:      am.BannerSN,
		BannerUrl:     am.BannerUrl,
		UserSN:        am.UserSN,
		Username:      am.Username,
		UserAvatar:    am.UserAvatar,
	}).Error
	if err != nil {
		tx.Rollback() // 如果文章创建失败，回滚事务
		global.Log.Errorf("article  UploadArticles err:%s\n", err.Error())
		return "", err
	}
	// 自动创建不存在的标签
	for _, tagTitle := range tagsStrList {
		var tagCount int64
		// 检查 tag 是否存在
		if err := tx.Table("tag_models").Where("title = ?", tagTitle).Count(&tagCount).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("check tag existence err:%s\n", err.Error())
			return "", err
		}
		// 如果不存在，创建 tag
		if tagCount == 0 {
			if err := tx.Create(&model.TagModel{Title: tagTitle}).Error; err != nil {
				tx.Rollback()
				global.Log.Errorf("auto create tag err:%s\n", err.Error())
				return "", err
			}
			global.Log.Infof("Auto created tag: %s", tagTitle)
		}
	}

	// 插入 ArticleTagModel 表
	for i := range tagsStrList {
		err = tx.Table("article_tag_models").
			Create(&model.ArticleTagModel{
				ArticleSN:    am.SN,
				ArticleTitle: am.Title,
				TagTitle:     tagsStrList[i],
			}).Error
		if err != nil {
			tx.Rollback() // 如果插入 ArticleTagModel 失败，回滚事务
			global.Log.Errorf("article  insert into article_tag_models err:%s\n", err.Error())
			return "", err
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		global.Log.Errorf("article  commit transaction err:%s\n", err.Error())
		return "", err
	}

	return code.StrCreateSucceed, nil
}

func PostArticleCollect(uSN int64, articleSN int64) (string, error) {
	tx := global.DB.Begin()

	// 1. 检查是否已经收藏
	var count int64
	if err := tx.Table("user_collect_models").
		Where("article_sn = ? and user_sn = ?", articleSN, uSN).Count(&count).Error; err != nil {
		global.Log.Errorf("user_collect_models count err: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	if count > 0 {
		// 已收藏，执行取消收藏逻辑
		if err := tx.Table("user_collect_models").
			Where("article_sn = ? and user_sn = ?", articleSN, uSN).
			Delete(&model.UserCollectModel{}).Error; err != nil {
			global.Log.Errorf("user_collect_models delete err: %s\n", err.Error())
			tx.Rollback()
			return "", err
		}

		// 更新收藏计数 -1
		if err := tx.Table("article_models").
			Where("sn = ?", articleSN).
			UpdateColumn("collects_count", gorm.Expr("collects_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("article_models update collects_count -1 err: %s\n", err.Error())
			return "", err
		}

		if err := tx.Commit().Error; err != nil {
			global.Log.Errorf("tx.Commit() error: %s\n", err.Error())
			tx.Rollback()
			return "", err
		}
		return "取消收藏成功", nil
	}

	// 2. 未收藏，执行创建收藏逻辑
	if err := tx.Create(&model.UserCollectModel{
		UserSN:    uSN,
		ArticleSN: articleSN,
	}).Error; err != nil {
		global.Log.Errorf("UserCollectModel create err: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	// 3. 更新收藏计数 +1
	if err := tx.Table("article_models").
		Where("sn = ?", articleSN).
		UpdateColumn("collects_count", gorm.Expr("collects_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("article_models update collects_count +1 err: %s\n", err.Error())
		return "", err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("tx.Commit() error: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	return "收藏成功", nil
}
