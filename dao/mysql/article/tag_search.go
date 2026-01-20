package article

import (
	dao_es "github.com/bramble555/blog/dao/es"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

// GetArticlesByTag 根据 Tag 获取文章列表 (ES -> MySQL 降级)
func GetArticlesByTag(tag string, pl *model.ParamList) (*model.ResponseArticleList, error) {
	var count int64
	var articles []model.ResponseArticle
	var esSNs []int64
	var err error
	useES := global.ES != nil

	// 1. 尝试使用 ES 搜索
	if useES {
		esSNs, err = dao_es.GetArticleSNsByTag(tag)
		if err != nil {
			global.Log.Errorf("GetArticlesByTag ES search failed: %v, fallback to MySQL", err)
			useES = false
		} else if len(esSNs) == 0 {
			// ES 搜索成功但无结果
			return &model.ResponseArticleList{
				List:      []model.ResponseArticle{},
				Count:     0,
				Page:      pl.Page,
				PageSize:  pl.Size,
				TotalPage: 0,
			}, nil
		}
	}

	db := global.DB.Table("article_models")

	// 2. 构建查询
	if useES {
		// ES 模式：直接用 SN 过滤
		db = db.Where("sn IN ?", esSNs)
		if err := db.Count(&count).Error; err != nil {
			return nil, err
		}
	} else {
		// MySQL 降级模式：使用 FIND_IN_SET 进行精确匹配 (避免 'Go' 匹配到 'Golang')
		db = db.Where("FIND_IN_SET(?, tags)", tag)
		if err := db.Count(&count).Error; err != nil {
			return nil, err
		}
	}

	// 3. 查询列表
	offset := (pl.Page - 1) * pl.Size
	selectCols := `article_models.sn, article_models.create_time, 
	article_models.update_time, article_models.title, 
	article_models.abstract, article_models.look_count, 
	article_models.comment_count, article_models.digg_count, 
	article_models.collects_count, article_models.tags, 
	article_models.banner_sn, article_models.banner_url, 
	article_models.user_sn, article_models.username, article_models.user_avatar`

	orderBy := pl.Order
	if useES && len(esSNs) > 0 {
		orderBy = buildFieldOrder(esSNs)
	}

	query := db.Select(selectCols).Order(orderBy).Limit(pl.Size).Offset(offset)

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	totalPage := 0
	if pl.Size > 0 {
		totalPage = int((count + int64(pl.Size) - 1) / int64(pl.Size))
	}

	return &model.ResponseArticleList{
		List:      articles,
		Count:     count,
		Page:      pl.Page,
		PageSize:  pl.Size,
		TotalPage: totalPage,
	}, nil
}
