package article

import (
	"strings"
	"time"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func IDListExist(pdl *model.ParamDeleteList) (bool, error) {
	var count int64
	err := global.DB.Table("article_models").Where("id IN ?", pdl.IDList).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error IDListExist: %v\n", err)
		return false, code.ErrorIDNotExit
	}
	return int(count) == len(pdl.IDList), nil
}
func IDExist(id uint) (bool, error) {
	var count int64
	err := global.DB.Table("article_models").Where("id = ?", id).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error IDExist: %v\n", err)
		return false, code.ErrorIDNotExit
	}
	return count > 0, nil
}
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

type MySQLArticleQueryService struct{}

func (q *MySQLArticleQueryService) GetArticlesListByParam(paq *model.ParamArticleQuery) (*[]model.ResponseArticle, error) {
	// 构造动态查询条件
	query := "SELECT id, title, content, tags FROM article_models WHERE 1=1"
	args := []interface{}{}

	// 如果传递了 title 参数，构造全文搜索条件
	if paq.Title != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+paq.Title+"%")
	}

	// 如果传递了 tags 参数
	if paq.Tags != "" {
		query += " AND tags LIKE ?"
		args = append(args, "%"+paq.Tags+"%")
	}

	// 如果传递了 content 参数
	if paq.Content != "" {
		query += " AND content LIKE ?"
		args = append(args, "%"+paq.Content+"%")
	}

	// 添加分页逻辑
	offset := (paq.Page - 1) * paq.Size
	query += " LIMIT ? OFFSET ?"
	args = append(args, paq.Size, offset)

	// 使用 gorm 的 Raw 方法执行 SQL 查询
	rows, err := global.DB.Table("article_models").Raw(query, args...).Rows()
	if err != nil {
		global.Log.Errorf("failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 解析结果
	var articles []model.ResponseArticle
	for rows.Next() {
		var article model.ResponseArticle
		var tags string
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &tags); err != nil {
			global.Log.Errorf("failed to scan row: %v", err)
			return nil, err
		}
		// 如果 tags 字段包含 `\n` 分割的多个值，进行拆分
		article.Tags = strings.Split(tags, "\n")
		articles = append(articles, article)
	}

	// 检查是否有错误
	if err := rows.Err(); err != nil {
		global.Log.Errorf("rows iteration error: %v", err)
		return nil, err
	}

	return &articles, nil
}

func GetArticlesDetail(id string) (*model.ArticleModel, error) {
	// 开始事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取文章详情
	am := &model.ArticleModel{}
	if err := tx.Table("article_models").
		Where("id = ?", id).
		First(am).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("select err:%s\n", err.Error())
		return nil, err
	}

	// 更新文章的浏览次数
	if err := tx.Table("article_models").
		Where("id = ?", id).
		UpdateColumn("look_count", am.LookCount+1).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("article_models Update err:%s\n", err.Error())
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("Transaction commit failed: %s\n", err.Error())
		return nil, err
	}

	return am, nil
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
func GetArticleCollect(uID uint) ([]model.ResponseArticle, error) {
	// 查询用户收藏的 article_id
	var articleIDList []uint
	if err := global.DB.Table("user_collect_models").
		Where("user_id = ?", uID).
		Select("article_id").
		Scan(&articleIDList).Error; err != nil {
		global.Log.Errorf("GetArticleCollect query err: %s\n", err.Error())
		return nil, err
	}

	// 根据 article_id 列表查询文章详情
	var articles []model.ResponseArticle
	if len(articleIDList) > 0 {
		if err := global.DB.Table("article_models").
			Where("id IN (?)", articleIDList).
			Find(&articles).Error; err != nil {
			global.Log.Errorf("GetArticleCollect find articles err: %s\n", err.Error())
			return nil, err
		}
	}
	return articles, nil
}

// GetUserCollectsCount 获取用户收藏的文章数量
func GetUserCollectsCount(uID uint, articleIDs []uint) (int64, error) {
	var count int64
	if err := global.DB.Table("user_collect_models").
		Where("user_id = ? AND article_id IN (?)", uID, articleIDs).
		Count(&count).Error; err != nil {
		global.Log.Errorf("GetUserCollectsCount err:%s\n", err.Error())
		return 0, err
	}
	return count, nil
}
func DeleteArticleCollect(uID uint, articleIDs []uint) (string, error) {
	tx := global.DB.Begin()

	// 封装回滚操作
	rollbackOnError := func(err error) error {
		tx.Rollback()
		return err
	}

	// 删除指定的文章收藏记录
	if err := tx.Table("user_collect_models").
		Where("user_id = ? AND article_id IN (?)", uID, articleIDs).
		Delete(&model.UserCollectModel{}).Error; err != nil {
		global.Log.Errorf("user_collect_models Delete err:%s\n", err.Error())
		return "", rollbackOnError(err)
	}

	// 更新每个文章的收藏数量
	for _, articleID := range articleIDs {
		// 查询当前文章的收藏数量
		var currentCount int64
		if err := tx.Table("article_models").
			Where("id = ?", articleID).
			Select("collects_count").
			Row().
			Scan(&currentCount); err != nil {
			global.Log.Errorf("article_models Select err:%s\n", err.Error())
			return "", rollbackOnError(err)
		}

		// 递减收藏数量
		newCount := currentCount - 1

		// 更新文章的收藏数量
		if err := tx.Table("article_models").
			Where("id = ?", articleID).
			Update("collects_count", newCount).Error; err != nil {
			global.Log.Errorf("article_models Update err:%s\n", err.Error())
			return "", rollbackOnError(err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return "", rollbackOnError(err)
	}

	return "取消收藏成功", nil
}
