package article

import (
	"time"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func GetArticlesList(pl *model.ParamList) (*[]model.ResponseArticle, error) {
	offset := (pl.Page - 1) * pl.Size
	res := []model.ResponseArticle{}
	err := global.DB.Table("article_models").
		Select("id, create_time, update_time, title, abstract, look_count, comment_count, digg_count, collects_count, category, tags, banner_id, banner_url, user_id, username, user_avatar").
		Order(pl.Order).
		Limit(pl.Size).
		Offset(offset).
		Find(&res).Error
	if err != nil {
		global.Log.Errorf("select err:%s\n", err.Error())
		return nil, err
	}
	return &res, nil
}
func GetArticlesDetail(id string) (*model.ArticleModel, error) {
	am := model.ArticleModel{}
	err := global.DB.Table("article_models").
		Where("id = ?", id).
		First(&am).Error

	if err != nil {
		// 错误处理，输出日志
		global.Log.Errorf("select err:%s\n", err.Error())
		return nil, err
	}
	return &am, nil
}
func GetArticlesCalendar() (*map[string]int, error) {
	// 想要包含今天，那就 + 1 天
	tomorrow := time.Now().AddDate(0, 0, 1)
	yearAgo := tomorrow.AddDate(-1, 0, 0)
	format := "2006-01-02"
	days := int(tomorrow.Sub(yearAgo).Hours() / 24)
	dateStrings := make(map[string]int, days)
	for i := 0; i < days; i++ {
		date := yearAgo.AddDate(0, 0, i) // 获取当前日期
		dateStr := date.Format(format)
		dateStrings[dateStr] = 0 // 格式化并加入到字符串切片
	}
	var result = []model.CalendarCount{}
	err := global.DB.Raw(`
		SELECT DATE_FORMAT(create_time, '%Y-%m-%d') as date, COUNT(id) as count
		FROM article_models
		WHERE create_time >= ? AND create_time <= ?
		GROUP BY date
	`, yearAgo, tomorrow).Scan(&result).Error
	for i := 0; i < len(result); i++ {
		if _, exists := dateStrings[result[i].Date]; exists {
			dateStrings[result[i].Date] = result[i].Count
		}
	}

	// 如果只想获取 有文章的日期，那就下面这个方法
	// err := global.DB.Table("article_models").
	// 	Select("DATE_FORMAT(create_time, '%Y-%m-%d') as date, count(id) as count").
	// 	Group("date").
	// 	Scan(&result).Error

	if err != nil {
		global.Log.Errorf("get calendar mysql err:%s\n", err.Error())
		return nil, err
	}
	return &dateStrings, nil
}
func GetArticlesTagsList(pl *model.ParamList) (*[]model.ResponseArticleTags, error) {
	offset := (pl.Page - 1) * pl.Size
	res := []model.ResponseArticleTags{}

	// 查询每个标签的文章数量以及文章标题列表
	err := global.DB.Table("article_tag_models").
		Select("tag_title, COUNT(article_title) AS count, GROUP_CONCAT(article_title ORDER BY article_title ASC) AS article_title_list, MIN(create_time) AS create_time").
		Group("tag_title").  // 按 tag_title 分组
		Order("count DESC"). // 根据请求的排序方式排序
		Limit(pl.Size).      // 限制返回的条目数
		Offset(offset).      // 分页偏移量
		Scan(&res).Error     // 执行查询并将结果扫描到 res 中

	if err != nil {
		global.Log.Errorf("select err:%s\n", err.Error())
		return nil, err
	}

	return &res, nil
}
func DeleteArticlesList(pdl *model.ParamDeleteList) (string, error) {
	// 删除 article_models 表中数据
	err := global.DB.Table("article_models").
		Where("id In ?", pdl.IDList).Delete(model.ArticleModel{}).Error
	if err != nil {
		global.Log.Errorf("delete article_models err:%s\n", err.Error())
		return "", err
	}

	// 删除 article_tag_models 表中数据
	err = global.DB.Table("article_tag_models").
		Where("article_id In ?", pdl.IDList).Delete(model.ArticleTagModel{}).Error
	if err != nil {
		global.Log.Errorf("delete article_tag_models err:%s\n", err.Error())
		return "", err
	}
	return "删除成功", nil
}
