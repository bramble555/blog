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
	now := time.Now()
	// 定义时间格式
	layout := "2006-01-02 15:04:05"
	// 将当前时间格式化为字符串
	tStr := now.Format(layout)
	err := global.DB.Create(&model.ArticleModel{
		// 创建一个 ArticleModel 实例并赋值
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
		global.Log.Errorf("article  UploadArticles err:%s\n", err.Error())
		return "", err
	}
	return code.CreateSucceed, nil
}
