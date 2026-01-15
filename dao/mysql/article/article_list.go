package article

import (
	"time"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"gorm.io/gorm"
)

func SNListExist(pdl *model.ParamDeleteList) (bool, error) {
	// 转换 SNList 为 []int64
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("SNListExist StringSliceToInt64Slice err: %s\n", err.Error())
		return false, err
	}

	var count int64
	err = global.DB.Table("article_models").Where("sn IN ?", snList).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error SNListExist: %v\n", err)
		return false, code.ErrorSNNotExit
	}
	return int(count) == len(snList), nil
}

// mysql 对应字段为 sn
func CheckSNExist(sn int64) (bool, error) {
	var count int64
	err := global.DB.Table("article_models").Where("sn = ?", sn).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error SNExist: %v\n", err)
		return false, code.ErrorSNNotExit
	}
	return count > 0, nil
}

func IsUserCollect(uSN int64, articleSN int64) (bool, error) {
	var count int64
	err := global.DB.Table("user_collect_models").
		Where("user_sn = ? AND article_sn = ?", uSN, articleSN).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// mysql 对应字段为 article_sn
// func SNExistByArticleSN(sn int64) (bool, error) {
// 	var count int64
// 	err := global.DB.Table("article_models").Where("sn = ?", sn).Count(&count).Error
// 	if err != nil {
// 		global.Log.Errorf("Error SNExist: %v\n", err)
// 		return false, code.ErrorSNNotExit
// 	}
// 	return count > 0, nil
// }

func GetArticlesList(pl *model.ParamList) (*[]model.ResponseArticle, error) {
	offset := (pl.Page - 1) * pl.Size
	res := []model.ResponseArticle{}
	err := global.DB.Table("article_models").
		Select("sn, create_time, update_time, title, abstract, look_count, comment_count, digg_count, collects_count, tags, banner_sn, banner_url, user_sn, username, user_avatar").
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

func GetArticlesListByParam(paq *model.ParamArticleQuery, uSN int64) (*model.ResponseArticleList, error) {
	articles := make([]model.ResponseArticle, 0)
	var count int64
	db := global.DB.Table("article_models")

	if paq.Title != "" {
		db = db.Where("title LIKE ?", "%"+paq.Title+"%")
	}
	if paq.Tags != "" {
		db = db.Where("tags LIKE ?", "%"+paq.Tags+"%")
	}
	if paq.Content != "" {
		db = db.Where("content LIKE ?", "%"+paq.Content+"%")
	}

	if err := db.Count(&count).Error; err != nil {
		global.Log.Errorf("GetArticlesListByParam count failed: %v", err)
		return nil, err
	}

	offset := (paq.Page - 1) * paq.Size
	// 默认按创建时间降序排序
	err := db.Order("create_time DESC").Limit(paq.Size).Offset(offset).Find(&articles).Error
	if err != nil {
		global.Log.Errorf("GetArticlesListByParam failed: %v", err)
		return nil, err
	}

	// 如果用户登录了，检查是否收藏
	if uSN != 0 && len(articles) > 0 {
		var articleSNs []int64
		for _, article := range articles {
			articleSNs = append(articleSNs, article.SN)
		}

		var collectSNs []int64
		global.DB.Table("user_collect_models").
			Where("user_sn = ? AND article_sn IN (?)", uSN, articleSNs).
			Pluck("article_sn", &collectSNs)

		collectMap := make(map[int64]bool)
		for _, sn := range collectSNs {
			collectMap[sn] = true
		}

		for i := range articles {
			if collectMap[articles[i].SN] {
				articles[i].IsCollect = true
			}
		}
	}

	return &model.ResponseArticleList{
		List:  articles,
		Count: count,
	}, nil
}

func GetArticlesDetail(sn string) (*model.ArticleModel, error) {
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
		Where("sn = ?", sn).
		First(am).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("article_models Query err:%s\n", err.Error())
		return nil, err
	}
	// 文章浏览量 +1
	err := tx.Table("article_models").Where("sn = ?", sn).
		UpdateColumn("look_count", gorm.Expr("look_count + ?", 1)).
		Error
	if err != nil {
		global.Log.Errorf("scan look_count err:%s\n", err.Error())
		tx.Rollback()
		return nil, err
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("Commit err:%s\n", err.Error())
		tx.Rollback()
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
		SELECT DATE_FORMAT(create_time, '%Y-%m-%d') as date, COUNT(sn) as count
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
	// 	Select("DATE_FORMAT(create_time, '%Y-%m-%d') as date, count(sn) as count").
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
	// 转换 SNList 为 []int64
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("DeleteArticlesList StringSliceToInt64Slice err: %s\n", err.Error())
		return "", err
	}

	// 删除 article_models 表中数据
	err = global.DB.Table("article_models").
		Where("sn In ?", snList).Delete(model.ArticleModel{}).Error
	if err != nil {
		global.Log.Errorf("delete article_models err:%s\n", err.Error())
		return "", err
	}

	// 删除 article_tag_models 表中数据
	err = global.DB.Table("article_tag_models").
		Where("article_sn In ?", snList).Delete(model.ArticleTagModel{}).Error
	if err != nil {
		global.Log.Errorf("delete article_tag_models err:%s\n", err.Error())
		return "", err
	}
	return "删除成功", nil
}
func GetArticleCollect(uSN int64) ([]model.ResponseArticle, error) {
	// 查询用户收藏的 article_sn
	var articleSNList []int64
	if err := global.DB.Table("user_collect_models").
		Where("user_sn = ?", uSN).
		Select("article_sn").
		Scan(&articleSNList).Error; err != nil {
		global.Log.Errorf("GetArticleCollect query err: %s\n", err.Error())
		return nil, err
	}

	// 根据 article_sn 列表查询文章详情
	var articles []model.ResponseArticle
	if len(articleSNList) > 0 {
		if err := global.DB.Table("article_models").
			Where("sn IN (?)", articleSNList).
			Find(&articles).Error; err != nil {
			global.Log.Errorf("GetArticleCollect find articles err: %s\n", err.Error())
			return nil, err
		}
		// 既然是收藏列表，IsCollect 肯定为 true
		for i := range articles {
			articles[i].IsCollect = true
		}
	}
	return articles, nil
}

// GetUserCollectsCount 获取用户收藏的文章数量
func GetUserCollectsCount(uSN int64, articleSNs []int64) (int64, error) {
	var count int64
	if err := global.DB.Table("user_collect_models").
		Where("user_sn = ? AND article_sn IN (?)", uSN, articleSNs).
		Count(&count).Error; err != nil {
		global.Log.Errorf("GetUserCollectsCount err:%s\n", err.Error())
		return 0, err
	}
	return count, nil
}
func DeleteArticleCollect(uSN int64, articleSNs []int64) (string, error) {
	tx := global.DB.Begin()

	// 封装回滚操作
	rollbackOnError := func(err error) error {
		tx.Rollback()
		return err
	}

	// 删除指定的文章收藏记录
	if err := tx.Table("user_collect_models").
		Where("user_sn = ? AND article_sn IN (?)", uSN, articleSNs).
		Delete(&model.UserCollectModel{}).Error; err != nil {
		global.Log.Errorf("user_collect_models Delete err:%s\n", err.Error())
		return "", rollbackOnError(err)
	}

	// 更新对应 article_models 中的 collects_count
	if err := tx.Table("article_models").
		Where("sn IN (?)", articleSNs).
		UpdateColumn("collects_count", gorm.Expr("collects_count - 1")).Error; err != nil {
		global.Log.Errorf("article_models Update collects_count err:%s\n", err.Error())
		return "", rollbackOnError(err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("Transaction commit failed: %s\n", err.Error())
		return "", err
	}

	return "取消收藏成功", nil
}
