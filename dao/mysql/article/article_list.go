package article

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	dao_es "github.com/bramble555/blog/dao/es"
	"github.com/bramble555/blog/dao/mysql/code"
	dao_redis "github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg/convert"
	"gorm.io/gorm"
)

// GetArticlesDetail 获取文章详情
func GetArticlesDetail(sn string) (*model.ArticleModel, error) {
	// 获取文章详情 (详情页不再直接在事务中对 MySQL 执行 +1 操作，极大提升性能)
	var am model.ArticleModel
	if err := global.DB.Table("article_models").Where("sn = ?", sn).First(&am).Error; err != nil {
		global.Log.Errorf("article_models Query err:%s\n", err.Error())
		return nil, err
	}

	// 1. 同时在 Redis 中异步增加浏览量 (大厂方案：Redis 优先写)
	go dao_redis.IncrArticleCount(am.SN, dao_redis.FieldLook, 1)

	// 2. 从 Redis 合并实时计数 (MySQL 存量 + Redis 增量)
	lookMap, _ := dao_redis.GetRedisArticlesCounts([]int64{am.SN}, dao_redis.FieldLook)
	if val, ok := lookMap[am.SN]; ok && val != 0 {
		am.LookCount += val
	}
	diggMap, _ := dao_redis.GetRedisArticlesCounts([]int64{am.SN}, dao_redis.FieldDigg)
	if val, ok := diggMap[am.SN]; ok && val != 0 {
		am.DiggCount += val
	}

	return &am, nil
}

// GetArticlesListByParam 获取文章列表 ,也可以通过 title tags content 进行搜索
func GetArticlesListByParam(paq *model.ParamArticleQuery, uSN int64) (*model.ResponseArticleList, error) {
	articles := make([]model.ResponseArticle, 0)
	var count int64

	// ===== 首页优化：优先从 Redis 获取 =====
	// 如果是首页请求（无搜索条件，第 1 页），直接从 Redis 获取最新文章 SN
	if paq.Title == "" && paq.Content == "" && paq.Page == 1 {
		sns, err := dao_redis.GetLatestArticleSNs(paq.Size)
		if err == nil && len(sns) > 0 {
			global.Log.Infof("Homepage: using Redis cache, got %d articles", len(sns))

			// 根据 SN 列表查询文章详情，使用 FIELD 保持顺序
			orderBy := buildFieldOrder(sns)
			selectCols := `article_models.sn, article_models.create_time, 
				article_models.update_time, article_models.title, 
				article_models.abstract, article_models.look_count, 
				article_models.comment_count, article_models.digg_count, 
				article_models.collects_count, article_models.tags, 
				article_models.banner_sn, article_models.banner_url, 
				article_models.user_sn, article_models.username, article_models.user_avatar`

			err := global.DB.Table("article_models").
				Where("article_models.sn IN (?)", sns).
				Select(selectCols).
				Order(orderBy).
				Find(&articles).Error

			if err != nil {
				global.Log.Warnf("Homepage: Redis cache hit but MySQL query failed: %v, fallback to normal logic", err)
				return nil, err
			} else {
				// 成功从 Redis 获取，合并实时计数并填充用户行为状态
				mergeRealtimeCounts(&articles) // 合并 Redis 实时计数 (大厂标准)

				if uSN != 0 && len(articles) > 0 {
					fillUserActions(&articles, uSN)
				}

				totalPage := 1
				if paq.Size > 0 {
					totalPage = (len(sns) + paq.Size - 1) / paq.Size
				}

				return &model.ResponseArticleList{
					List:      articles,
					Count:     int64(len(sns)),
					Page:      paq.Page,
					PageSize:  paq.Size,
					TotalPage: totalPage,
				}, nil
			}
		}
		// Redis 失败或无数据，降级到正常逻辑
		global.Log.Infof("Homepage: Redis cache miss, using MySQL")
	}

	// 如果有 title 或 content 搜索关键词，首先使用 ES 搜索

	db := global.DB.Table("article_models")
	useESSearch := false // 标记是否使用 ES 搜索
	var esSNs []int64    // ES 搜索返回的文章 SN 列表

	// ===== ES 搜索逻辑 =====
	searchKeyword := ""
	if paq.Title != "" {
		searchKeyword = paq.Title
	} else if paq.Content != "" {
		searchKeyword = paq.Content
	}

	if searchKeyword != "" && global.ES != nil {
		// 使用 ES 搜索
		useESSearch = true
		// 获取匹配的文章 SN 列表
		var err error
		esSNs, err = dao_es.GetArticleSNsByKeyword(searchKeyword)
		if err != nil {
			global.Log.Errorf("GetArticlesListByParam ES search failed: %v, fallback to MySQL", err)
			// ES 搜索失败，降级到 MySQL LIKE 查询
			useESSearch = false
		} else if len(esSNs) == 0 {
			// ES 搜索无结果，直接返回空列表
			return &model.ResponseArticleList{
				List:      []model.ResponseArticle{},
				Count:     0,
				Page:      paq.Page,
				PageSize:  paq.Size,
				TotalPage: 0,
			}, nil
		}

		// 使用 ES 返回的 SN 列表过滤
		db = db.Where("article_models.sn IN ?", esSNs)
	}

	// 当 ES 不可用的时候,用 mysql 搜索
	if !useESSearch {
		if paq.Title != "" {
			db = db.Where("article_models.title LIKE ?", "%"+paq.Title+"%")
		}
		if paq.Content != "" {
			db = db.Where("article_models.content LIKE ?", "%"+paq.Content+"%")
		}
	}

	if err := db.Count(&count).Error; err != nil {
		global.Log.Errorf("GetArticlesListByParam count failed: %v", err)
		return nil, err
	}

	offset := (paq.Page - 1) * paq.Size
	query := db
	selectCols := `article_models.sn, article_models.create_time, 
	article_models.update_time, article_models.title, 
	article_models.abstract, article_models.look_count, 
	article_models.comment_count, article_models.digg_count, 
	article_models.collects_count, article_models.tags, 
	article_models.banner_sn, article_models.banner_url, 
	article_models.user_sn, article_models.username, article_models.user_avatar`

	// 如果使用 ES 搜索，按照 ES 返回的顺序排序（相关度）
	// 否则按创建时间降序排序
	orderBy := "article_models.create_time DESC"
	if useESSearch && len(esSNs) > 0 {
		// 使用 FIELD() 函数保持 ES 的相关度排序
		// 注意：这里需要构建 FIELD(sn, id1, id2, ...) 语句
		orderBy = buildFieldOrder(esSNs)
	}

	err := query.Select(selectCols).Order(orderBy).
		Limit(paq.Size).
		Offset(offset).
		Find(&articles).Error
	if err != nil {
		global.Log.Errorf("GetArticlesListByParam failed: %v", err)
		return nil, err
	}

	// 查询成功后，合并 Redis 中实时的浏览、点赞、收藏计数器
	if len(articles) > 0 {
		mergeRealtimeCounts(&articles)
	}

	// 如果用户登录了，检查是否收藏和点赞
	if uSN != 0 && len(articles) > 0 {
		fillUserActions(&articles, uSN)
	}
	totalPage := 0
	if paq.Size > 0 {
		totalPage = int((count + int64(paq.Size) - 1) / int64(paq.Size))
	}
	return &model.ResponseArticleList{
		List:      articles,
		Count:     count,
		Page:      paq.Page,
		PageSize:  paq.Size,
		TotalPage: totalPage,
	}, nil
}

// CheckSnListExist 检查文章 SNList 是否存在
func SNListExist(pdl *model.ParamDeleteList) (bool, error) {
	// 转换 SNList 为 []int64
	snList, err := convert.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("SNListExist StringSliceToInt64Slice err: %s\n", err.Error())
		return false, err
	}
	// 查询文章数量
	var count int64
	err = global.DB.Table("article_models").Where("sn IN ?", snList).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error SNListExist: %v\n", err)
		return false, code.ErrorSNNotExist
	}
	return int(count) == len(snList), nil
}

// CheckArticleExist 检查文章是否存在
func CheckSNExist(sn int64) (bool, error) {
	var count int64
	err := global.DB.Table("article_models").Where("sn = ?", sn).Count(&count).Error
	if err != nil {
		global.Log.Errorf("Error SNExist: %v\n", err)
		return false, code.ErrorSNNotExist
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
	snList, err := convert.StringSliceToInt64Slice(pdl.SNList)
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

	// 从 Redis 缓存中删除文章（首页最新文章列表）
	// 不阻塞主流程，Redis 删除失败只记录日志
	for _, sn := range snList {
		if err := dao_redis.RemoveArticleFromLatest(sn); err != nil {
			global.Log.Warnf("Failed to remove article %d from Redis cache: %v", sn, err)
		}
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

// buildFieldOrder 构建 FIELD() 函数用于保持 ES 的相关度排序
func buildFieldOrder(sns []int64) string {
	if len(sns) == 0 {
		return "article_models.create_time DESC"
	}

	// 构建 FIELD(article_models.sn, id1, id2, ...) 语句
	var snStrs []string
	for _, sn := range sns {
		snStrs = append(snStrs, strconv.FormatInt(sn, 10))
	}
	return fmt.Sprintf("FIELD(article_models.sn, %s)", strings.Join(snStrs, ","))
}

// mergeRealtimeCounts 从 Redis 批量获取计数器实时数据并合并到文章模型中
func mergeRealtimeCounts(articles *[]model.ResponseArticle) {
	if len(*articles) == 0 {
		return
	}

	sns := make([]int64, len(*articles))
	for i, a := range *articles {
		sns[i] = a.SN
	}

	// 批量获取各类计数
	lookCounts, _ := dao_redis.GetRedisArticlesCounts(sns, dao_redis.FieldLook)
	diggCounts, _ := dao_redis.GetRedisArticlesCounts(sns, dao_redis.FieldDigg)

	// 合并 (增量合并逻辑：MySQL 存量 + Redis 增量)
	for i := range *articles {
		sn := (*articles)[i].SN
		if val, ok := lookCounts[sn]; ok && val != 0 {
			(*articles)[i].LookCount += val
		}
		if val, ok := diggCounts[sn]; ok && val != 0 {
			(*articles)[i].DiggCount += val
		}
	}
}

// fillUserActions 填充用户的点赞和收藏状态 (全量走 MySQL)
func fillUserActions(articles *[]model.ResponseArticle, uSN int64) {
	if len(*articles) == 0 {
		return
	}

	// 提取文章 SN 列表
	articleSNs := make([]int64, 0, len(*articles))
	for _, article := range *articles {
		articleSNs = append(articleSNs, article.SN)
	}

	// 使用 MySQL 批量查询用户点赞情况
	diggMap := fillDiggFromMySQL(uSN, articleSNs)

	// 使用 MySQL 批量查询用户收藏情况
	collectMap := fillCollectFromMySQL(uSN, articleSNs)

	// 填充到文章列表
	for i := range *articles {
		(*articles)[i].IsDigg = diggMap[(*articles)[i].SN]
		(*articles)[i].IsCollect = collectMap[(*articles)[i].SN]
	}
}

// fillDiggFromMySQL 从 MySQL 查询用户点赞状态（降级方案）
func fillDiggFromMySQL(uSN int64, articleSNs []int64) map[int64]bool {
	var diggSNs []int64
	err := global.DB.Table("user_digg_models").
		Where("user_sn = ? AND article_sn IN (?)", uSN, articleSNs).
		Pluck("article_sn", &diggSNs).Error
	if err != nil {
		global.Log.Errorf("fillDiggFromMySQL: query failed: %v", err)
		return make(map[int64]bool)
	}

	diggMap := make(map[int64]bool, len(diggSNs))
	for _, sn := range diggSNs {
		diggMap[sn] = true
	}
	return diggMap
}

// fillCollectFromMySQL 从 MySQL 查询用户收藏状态（降级方案）
func fillCollectFromMySQL(uSN int64, articleSNs []int64) map[int64]bool {
	var collectSNs []int64
	err := global.DB.Table("user_collect_models").
		Where("user_sn = ? AND article_sn IN (?)", uSN, articleSNs).
		Pluck("article_sn", &collectSNs).Error
	if err != nil {
		global.Log.Errorf("fillCollectFromMySQL: query failed: %v", err)
		return make(map[int64]bool)
	}

	collectMap := make(map[int64]bool, len(collectSNs))
	for _, sn := range collectSNs {
		collectMap[sn] = true
	}
	return collectMap
}
