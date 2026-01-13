package article

import (
	"time"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
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

	// 生成 SN
	am.SN = global.Snowflake.GetID()

	now := time.Now()
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
		Category:      am.Category,
		Source:        am.Source,
		Link:          am.Link,
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

	tagsStrList := []string(am.Tags)
	n := len(am.Tags)
	if n == 0 {
		tx.Commit() // 如果没有标签，直接提交事务
		return code.StrCreateSucceed, nil
	}
	// 插入 ArticleTagModel 表
	for i := range tagsStrList {
		err = tx.Table("article_tag_models").
			Create(&model.ArticleTagModel{
				ArticleSN:    am.SN, // 使用生成的 SN
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
