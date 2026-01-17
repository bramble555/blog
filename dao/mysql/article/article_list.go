package article

import (
	"time"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"gorm.io/gorm"
)

// CheckSnListExist 检查文章 SNList 是否存在
func SNListExist(pdl *model.ParamDeleteList) (bool, error) {
	// 转换 SNList 为 []int64
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("SNListExist StringSliceToInt64Slice err: %s\n", err.Error())
		return false, err
	}
	// 查询文章数量
	var count int64
	err = global.DB.Table("article_models").Where("sn IN ?", snList).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error SNListExist: %v\n", err)
		return false, code.ErrorSNNotExit
	}
	return int(count) == len(snList), nil
}

// CheckArticleExist 检查文章是否存在
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

// GetArticlesListByParam 获取文章列表 ,也可以通过 title tags content 进行搜索
func GetArticlesListByParam(paq *model.ParamArticleQuery, uSN int64) (*model.ResponseArticleList, error) {
	articles := make([]model.ResponseArticle, 0)
	var count int64
	db := global.DB.Table("article_models")

	hasTagFilter := false

	if paq.Title != "" {
		db = db.Where("article_models.title LIKE ?", "%"+paq.Title+"%")
	}
	if paq.Tags != "" {
		tagList := pkg.ParseTagsStringSlice(paq.Tags)
		if len(tagList) > 0 {
			hasTagFilter = true
			db = db.Joins("JOIN article_tag_models atm ON atm.article_sn = article_models.sn").
				Where("atm.tag_title IN ?", tagList)
		}
	}
	if paq.Content != "" {
		db = db.Where("article_models.content LIKE ?", "%"+paq.Content+"%")
	}

	if hasTagFilter {
		if err := db.Distinct("article_models.sn").Count(&count).Error; err != nil {
			global.Log.Errorf("GetArticlesListByParam count failed: %v", err)
			return nil, err
		}
	} else if err := db.Count(&count).Error; err != nil {
		global.Log.Errorf("GetArticlesListByParam count failed: %v", err)
		return nil, err
	}

	offset := (paq.Page - 1) * paq.Size
	// 默认按创建时间降序排序
	query := db
	if hasTagFilter {
		query = query.Distinct("article_models.sn")
	}
	err := query.Order("article_models.create_time DESC").
		Limit(paq.Size).
		Offset(offset).
		Find(&articles).Error
	if err != nil {
		global.Log.Errorf("GetArticlesListByParam failed: %v", err)
		return nil, err
	}

	// 如果用户登录了，检查是否收藏和点赞
	if uSN != 0 && len(articles) > 0 {
		var articleSNs []int64
		for _, article := range articles {
			articleSNs = append(articleSNs, article.SN)
		}

		// 查询用户收藏的文章
		var collectSNs []int64
		err := global.DB.Table("user_collect_models").
			Where("user_sn = ? AND article_sn IN (?)", uSN, articleSNs).
			Pluck("article_sn", &collectSNs).Error
		if err != nil {
			global.Log.Errorf("GetArticlesCollectListByParam pluck failed: %v", err)
			return nil, err
		}

		// 如果有收藏的文章，设置 IsCollect 为 true
		// 使用 map, 避免 O(n^2)复杂度
		collectMap := make(map[int64]bool, len(collectSNs)) // 给 map 预设容量，性能更好
		for _, sn := range collectSNs {
			collectMap[sn] = true
		}

		for i := range articles {
			// 这样即便 Map 是空的，逻辑也成立（返回 false）
			articles[i].IsCollect = collectMap[articles[i].SN]
		}

		// 查询用户是否有点赞的文章
		var diggSNs []int64
		err = global.DB.Table("user_digg_models").
			Where("user_sn = ? AND article_sn IN (?)", uSN, articleSNs).
			Pluck("article_sn", &diggSNs).Error
		if err != nil {
			global.Log.Errorf("GetArticlesDiggListByParam pluck failed: %v", err)
			return nil, err
		}

		// 如果有点赞的文章, 设置 IsDigg 为 true
		diggMap := make(map[int64]bool, len(diggSNs))
		for _, sn := range diggSNs {
			diggMap[sn] = true
		}
		for i := range articles {
			articles[i].IsDigg = diggMap[articles[i].SN]
		}

		// 最后统一一次遍历，填充 IsCollect 和 IsDigg
		for i := range articles {
			articles[i].IsCollect = collectMap[articles[i].SN]
			articles[i].IsDigg = diggMap[articles[i].SN]
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

// GetArticlesCalendar 获取文章发布日历
func GetArticlesCalendar() (map[string]int, error) {
	// 1. 确定时间范围
	// 统一以 "天" 为单位处理，忽略时分秒带来的边缘计算复杂性
	now := time.Now()
	// 获取当天的 00:00:00 (作为基准)
	todayStart := time.Date(now.Year(), now.Month(),
		now.Day(), 0, 0, 0, 0, now.Location())

	// 结束时间: 今天的 23:59:59 (用于 SQL 过滤)
	endDate := todayStart.AddDate(0, 0, 1).Add(-time.Nanosecond)
	// 开始时间: 365天前 (用于 SQL 过滤)
	startDate := todayStart.AddDate(-1, 0, 0)

	// 2. 初始化 Map 并预填充 0
	// 预分配容量，避免 map 在扩容时重新哈希，366 涵盖闰年
	dateMap := make(map[string]int, 366)

	// 使用 AddDate 循环，
	for d := startDate; !d.After(todayStart); d = d.AddDate(0, 0, 1) {
		dateMap[d.Format("2006-01-02")] = 0
	}

	// 3. 执行 SQL 查询
	var results []model.CalendarCount

	// date 字段直接映射到 model，无需再手动 Scan 进特定变量
	err := global.DB.Model(&model.ArticleModel{}).
		Select("DATE_FORMAT(create_time, '%Y-%m-%d') as date, COUNT(*) as count").
		Where("create_time BETWEEN ? AND ?", startDate, endDate).
		Group("date").
		Scan(&results).Error

	if err != nil {
		global.Log.Errorf("GetArticlesCalendar db query err: %s", err.Error())
		return nil, err
	}

	// 4. 合并数据
	for _, r := range results {
		// 发表文章多的天数,颜色越深,交给前端处理
		dateMap[r.Date] = r.Count
	}

	// 直接返回 map，不需要指针
	return dateMap, nil
}

// DeleteArticlesList 删除文章列表
func DeleteArticlesList(pdl *model.ParamDeleteList) (string, error) {
	// 转换 SNList 为 []int64
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("DeleteArticlesList StringSliceToInt64Slice err: %s\n", err.Error())
		return "", err
	}
	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 避免下面重复写回滚
	rollbackOnError := func(err error) (string, error) {
		tx.Rollback()
		return "", err
	}

	var commentSNList []int64
	if err := tx.Table("comment_models").
		Where("article_sn IN ?", snList).
		Pluck("sn", &commentSNList).Error; err != nil {
		global.Log.Errorf("DeleteArticlesList pluck comment sn err:%s\n", err.Error())
		return rollbackOnError(err)
	}
	// 删除 user_comment_digg_models 表中数据
	if len(commentSNList) > 0 {
		if err := tx.Table("user_comment_digg_models").
			Where("comment_sn IN ?", commentSNList).
			Delete(&model.UserCommentDiggModel{}).Error; err != nil {
			global.Log.Errorf("DeleteArticlesList delete user_comment_digg_models err:%s\n", err.Error())
			return rollbackOnError(err)
		}
	}
	// 删除 comment_models 表中数据
	if err := tx.Table("comment_models").
		Where("article_sn IN ?", snList).
		Delete(&model.CommentModel{}).Error; err != nil {
		global.Log.Errorf("DeleteArticlesList delete comment_models err:%s\n", err.Error())
		return rollbackOnError(err)
	}

	// 删除 user_digg_models 表中数据
	if err := tx.Table("user_digg_models").
		Where("article_sn IN ?", snList).
		Delete(&model.UserDiggModel{}).Error; err != nil {
		global.Log.Errorf("DeleteArticlesList delete user_digg_models err:%s\n", err.Error())
		return rollbackOnError(err)
	}

	// 删除 user_collect_models 表中数据
	if err := tx.Table("user_collect_models").
		Where("article_sn IN ?", snList).
		Delete(&model.UserCollectModel{}).Error; err != nil {
		global.Log.Errorf("DeleteArticlesList delete user_collect_models err:%s\n", err.Error())
		return rollbackOnError(err)
	}

	// 删除 article_tag_models 表中数据
	if err := tx.Table("article_tag_models").
		Where("article_sn IN ?", snList).
		Delete(&model.ArticleTagModel{}).Error; err != nil {
		global.Log.Errorf("DeleteArticlesList delete article_tag_models err:%s\n", err.Error())
		return rollbackOnError(err)
	}

	// 删除 article_models 表中数据
	if err := tx.Table("article_models").
		Where("sn IN ?", snList).
		Delete(&model.ArticleModel{}).Error; err != nil {
		global.Log.Errorf("DeleteArticlesList delete article_models err:%s\n", err.Error())
		return rollbackOnError(err)
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("DeleteArticlesList commit err:%s\n", err.Error())
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
