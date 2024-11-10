package article

import (
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

// UpdateArticles 更新文章，如果传入 tags，也需要传入 title
func UpdateArticles(id uint, uf map[string]any) (string, error) {
	tx := global.DB.Begin()

	// 处理 content 字段
	if content, exists := uf["content"]; exists {
		uf["content"] = content
		// 如果有 title 需要更新 abstract 字段
		c := []rune(content.(string)) // 断言为字符串
		uf["abstract"] = string(c)    // 如果 content 长度小于 100，则直接取 content
		if len(c) > 100 {
			uf["abstract"] = string(c[:100]) // 提取前100个字符作为简介
		}
	}

	// 更新 article_models 表
	if err := tx.Table("article_models").Where("id = ?", id).Updates(uf).Error; err != nil {
		global.Log.Errorf("article_models update err: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	// 如果存在 title，更新 article_tag_models 表中的 article_title
	if title, exists := uf["title"]; exists {
		if err := tx.Table("article_tag_models").Where("article_id = ?", id).
			Updates(map[string]interface{}{"article_title": title}).Error; err != nil {
			global.Log.Errorf("article_tag_models update err: %s\n", err.Error())
			tx.Rollback() // 回滚事务
			return "", err
		}
	}

	// 如果存在 tags，需要处理标签更新
	if tagsTitle, exists := uf["tags"]; exists {
		// 删除原来的 tags
		if err := tx.Table("article_tag_models").Where("article_id = ?", id).
			Delete(model.ArticleTagModel{}).Error; err != nil {
			global.Log.Errorf("article_tag_models delete err: %s\n", err.Error())
			tx.Rollback() // 回滚事务
			return "", err
		}

		// 处理 tags: 检查每个标签是否存在并插入新的 tags
		if tags, ok := tagsTitle.([]interface{}); ok {
			var articleTagModels []model.ArticleTagModel
			for _, tag := range tags {
				t, ok := tag.(string)
				if !ok {
					global.Log.Errorf("tags 断言失败\n")
					tx.Rollback() // 回滚事务
					return "", code.ErrorAssertionFailed
				}

				// 检查 tag 是否存在
				var count int64
				if err := global.DB.Table("tag_models").Where("title = ?", t).Count(&count).Error; err != nil {
					global.Log.Errorf("tag check failed: %s\n", err.Error())
					tx.Rollback()
					return "", err
				}
				if count != 1 {
					global.Log.Errorf("tag:%s不存在", t)
					tx.Rollback() // 回滚事务
					return "", code.ErrorTagNotExit
				}

				// 添加标签
				articleTagModels = append(articleTagModels, model.ArticleTagModel{
					ArticleID: id,
					TagTitle:  t,
				})
			}

			// 批量插入标签
			if err := tx.Table("article_tag_models").CreateInBatches(articleTagModels, 100).Error; err != nil {
				global.Log.Errorf("article_tag_models batch insert err: %s\n", err.Error())
				tx.Rollback() // 回滚事务
				return "", err
			}
		} else {
			global.Log.Errorf("tags 类型断言失败\n")
			tx.Rollback()
			return "", code.ErrorAssertionFailed
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("tx.Commit() error: %s\n", err.Error())
		tx.Rollback() // 回滚事务
		return "", err
	}

	return code.StrUpdateSucceed, nil
}
func PostArticleCollect(uID uint, articleID uint) (string, error) {
	tx := global.DB.Begin()

	// 检查是否已经收藏
	var count int64
	if err := tx.Table("user_collect_models").
		Where("article_id = ? and user_id = ?", articleID, uID).Count(&count).Error; err != nil {
		global.Log.Errorf("user_collect_models count err: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}
	if count > 0 {
		return "已收藏", nil
	}

	// 创建收藏记录
	if err := tx.Create(&model.UserCollectModel{
		UserID:    uID,
		ArticleID: articleID,
	}).Error; err != nil {
		global.Log.Errorf("UserCollectModel create err: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	// 查询当前的 collects_count
	var currentCount int64
	if err := tx.Table("article_models").
		Select("collects_count").
		Where("id = ?", articleID).
		Scan(&currentCount).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("article_models select err: %s\n", err.Error())
		return "", err
	}

	// 原子更新收藏计数
	if err := tx.Table("article_models").
		Where("id = ?", articleID).
		Update("collects_count", currentCount+1).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("article_models update err: %s\n", err.Error())
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
