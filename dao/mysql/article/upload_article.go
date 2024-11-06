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
		// 记录错误日志
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

	now := time.Now()
	// 定义时间格式
	layout := "2006-01-02 15:04:05"
	// 将当前时间格式化为字符串
	tStr := now.Format(layout)

	// 创建一个 ArticleModel 实例并赋值
	err := tx.Create(&model.ArticleModel{
		CreateTime:    tStr,
		UpdateTime:    tStr,
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
		BannerID:      am.BannerID,
		BannerUrl:     am.BannerUrl,
		UserID:        am.UserID,
		Username:      am.Username,
		UserAvatar:    am.UserAvatar,
	}).Error
	if err != nil {
		tx.Rollback() // 如果文章创建失败，回滚事务
		global.Log.Errorf("article  UploadArticles err:%s\n", err.Error())
		return "", err
	}
	// 获取标签并查询标签 ID
	tagsStrList := []string(am.Tags)
	n := len(am.Tags)
	idList := make([]uint, n)
	if n == 0 {
		tx.Commit() // 如果没有标签，直接提交事务
		return code.CreateSucceed, nil
	}
	err = tx.Table("tag_models").Select("id").Where("title IN (?)", tagsStrList).
		Scan(&idList).Error
	if err != nil {
		tx.Rollback() // 如果查询标签 ID 失败，回滚事务
		global.Log.Errorf("article  select tags id err:%s\n", err.Error())
		return "", err
	}
	global.Log.Debugf("idlist:%+v", idList)
	// 插入 ArticleTagModel 表
	for _, id := range idList {
		err = tx.Table("article_tag_models").
			Create(&model.ArticleTagModel{
				ArticleID: am.ID,
				TagID:     id,
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

	return code.CreateSucceed, nil
}