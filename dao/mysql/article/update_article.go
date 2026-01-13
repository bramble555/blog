package article

import (
	"strconv"
	"strings"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"gorm.io/gorm"
)

// UpdateArticles 更新文章，如果传入 tags，也需要传入 title
func UpdateArticles(sn int64, uf map[string]any) (string, error) {
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

	// 处理 tags 字段，防止 GORM 无法识别 []interface{}
	var tagsToUpdate []string
	var hasTags bool
	if tags, exists := uf["tags"]; exists {
		if tList, ok := tags.([]interface{}); ok {
			hasTags = true
			for _, tag := range tList {
				if t, ok := tag.(string); ok {
					tagsToUpdate = append(tagsToUpdate, t)
				}
			}
			// ctype.Array 使用 \n 分隔，更新 article_models 表中的 tags 字段
			uf["tags"] = strings.Join(tagsToUpdate, "\n")
		}
	}
	// 处理 banner_sn，防止 frontend 传字符串导致 GORM 错误
	if bannerSN, exists := uf["banner_sn"]; exists {
		if bStr, ok := bannerSN.(string); ok {
			if bSN, err := strconv.ParseInt(bStr, 10, 64); err == nil {
				uf["banner_sn"] = bSN
			}
		}
	}

	// 更新 article_models 表
	if err := tx.Table("article_models").Where("sn = ?", sn).Updates(uf).Error; err != nil {
		global.Log.Errorf("article_models update err: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	// 获取文章标题，用于更新关联表
	var articleTitle string
	if title, exists := uf["title"]; exists {
		articleTitle = title.(string)
		// 如果存在 title，同步更新 article_tag_models 表中的 article_title
		if err := tx.Table("article_tag_models").Where("article_sn = ?", sn).
			Updates(map[string]interface{}{"article_title": articleTitle}).Error; err != nil {
			global.Log.Errorf("article_tag_models update err: %s\n", err.Error())
			tx.Rollback() // 回滚事务
			return "", err
		}
	} else {
		// 如果 uf 中没有 title，从数据库获取当前的 title
		if err := tx.Table("article_models").Select("title").Where("sn = ?", sn).Scan(&articleTitle).Error; err != nil {
			global.Log.Errorf("get article title err: %s\n", err.Error())
			tx.Rollback()
			return "", err
		}
	}

	// 如果存在 tags，需要处理标签更新
	if hasTags {
		// 去重 tagsToUpdate
		uniqueTags := make([]string, 0, len(tagsToUpdate))
		tagMap := make(map[string]struct{})
		for _, t := range tagsToUpdate {
			if _, ok := tagMap[t]; !ok {
				tagMap[t] = struct{}{}
				uniqueTags = append(uniqueTags, t)
			}
		}

		// 删除原来的 tags
		if err := tx.Table("article_tag_models").Where("article_sn = ?", sn).
			Delete(&model.ArticleTagModel{}).Error; err != nil {
			global.Log.Errorf("article_tag_models delete err: %s\n", err.Error())
			tx.Rollback() // 回滚事务
			return "", err
		}

		// 处理 tags: 检查每个标签是否存在并插入新的 tags
		var articleTagModels []model.ArticleTagModel
		for _, tag := range uniqueTags {
			// 检查 tag 是否存在
			var count int64
			if err := tx.Table("tag_models").Where("title = ?", tag).Count(&count).Error; err != nil {
				global.Log.Errorf("tag check failed: %s\n", err.Error())
				tx.Rollback()
				return "", err
			}
			if count == 0 {
				global.Log.Errorf("tag:%s不存在", tag)
				tx.Rollback() // 回滚事务
				return "", code.ErrorTagNotExit
			}

			// 添加标签
			articleTagModels = append(articleTagModels, model.ArticleTagModel{
				ArticleSN:    sn,
				TagTitle:     tag,
				ArticleTitle: articleTitle,
			})
		}

		// 批量插入标签
		if len(articleTagModels) > 0 {
			if err := tx.Table("article_tag_models").CreateInBatches(articleTagModels, 100).Error; err != nil {
				global.Log.Errorf("article_tag_models batch insert err: %s\n", err.Error())
				tx.Rollback() // 回滚事务
				return "", err
			}
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
func PostArticleCollect(uSN int64, articleSN int64) (string, error) {
	tx := global.DB.Begin()

	// 1. 检查是否已经收藏
	var count int64
	if err := tx.Table("user_collect_models").
		Where("article_sn = ? and user_sn = ?", articleSN, uSN).Count(&count).Error; err != nil {
		global.Log.Errorf("user_collect_models count err: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	if count > 0 {
		// 已收藏，执行取消收藏逻辑
		if err := tx.Table("user_collect_models").
			Where("article_sn = ? and user_sn = ?", articleSN, uSN).
			Delete(&model.UserCollectModel{}).Error; err != nil {
			global.Log.Errorf("user_collect_models delete err: %s\n", err.Error())
			tx.Rollback()
			return "", err
		}

		// 原子更新收藏计数 -1
		if err := tx.Table("article_models").
			Where("sn = ?", articleSN).
			UpdateColumn("collects_count", gorm.Expr("collects_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("article_models update collects_count -1 err: %s\n", err.Error())
			return "", err
		}

		if err := tx.Commit().Error; err != nil {
			global.Log.Errorf("tx.Commit() error: %s\n", err.Error())
			tx.Rollback()
			return "", err
		}
		return "取消收藏成功", nil
	}

	// 2. 未收藏，执行创建收藏逻辑
	if err := tx.Create(&model.UserCollectModel{
		UserSN:    uSN,
		ArticleSN: articleSN,
	}).Error; err != nil {
		global.Log.Errorf("UserCollectModel create err: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	// 3. 原子更新收藏计数 +1
	if err := tx.Table("article_models").
		Where("sn = ?", articleSN).
		UpdateColumn("collects_count", gorm.Expr("collects_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("article_models update collects_count +1 err: %s\n", err.Error())
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
