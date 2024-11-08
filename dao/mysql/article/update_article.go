package article

import (
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

// UpdateArticles 更新文章，如果传入 tags，也需要传入 title
func UpdateArticles(id uint, uf model.UpdatedFields) (string, error) {
	tx := global.DB.Begin()
	updates := make(map[string]interface{})
	// 每个字段都需要断言
	// 处理 content 字段
	if content, exists := (uf)["content"]; exists {
		updates["content"] = content
		// 如果有 title 需要更新 abstract 字段
		c := []rune(content.(string)) // 断言为字符串
		var abstract string
		if len(c) >= 100 {
			abstract = string(c[:100]) // 提取前100个字符作为简介
		} else {
			abstract = string(c)
		}
		updates["abstract"] = abstract
	}
	// 处理 title 字段
	if title, exists := (uf)["title"]; exists {
		updates["title"] = title
	}
	// 更新 article_models 表
	err := tx.Table("article_models").Where("id = ?", id).Updates(updates).Error
	if err != nil {
		global.Log.Errorf("article_models update err:%s\n", err.Error())
		tx.Rollback()
		return "", err
	}
	var articleTitle string
	// 查看是否有 article_title
	// 无论上面是否有 title 这个字段，article_models 表中的 article_title 已经更新
	// 下面只需要更新 article_tag_models 表中的 article_title
	if title, exists := (uf)["title"]; exists {
		err := tx.Table("article_tag_models").Where("article_id = ?", id).
			Updates(map[string]interface{}{
				"article_title": title,
			}).Error
		if err != nil {
			global.Log.Errorf("article_tag_models update err: %s\n", err.Error())
			tx.Rollback() // 回滚事务
			return "", err
		}
		var ok bool
		articleTitle, ok = title.(string)
		if !ok {
			global.Log.Errorf("articleTitle 断言失败\n")
			return "", code.ErrorAssertionFailed
		}
	}

	// 查看是否有 tags_title(有 tags_title 的前提是 必须有 article_title)
	if tagsTitle, exists := (uf)["tags"]; exists {
		// 有 tags_title 那就先删除原来的 tags
		err := tx.Table("article_tag_models").Where("article_id = ?", id).
			Delete(model.ArticleTagModel{}).Error
		if err != nil {
			global.Log.Errorf("article_tag_models delete err: %s\n", err.Error())
			tx.Rollback() // 回滚事务
			return "", err
		}
		tags, ok := tagsTitle.([]interface{})
		if !ok {
			global.Log.Errorf("tags 断言失败\n")
			return "", code.ErrorAssertionFailed
		}
		// 批量插入标签,优化 RTT
		var articleTagModels []model.ArticleTagModel
		for _, tag := range tags {
			t := tag.(string)
			// 检查 tag 是否存在
			var count int64
			global.DB.Table("tag_models").Where("title = ?", t).Count(&count)
			if count != 1 {
				global.Log.Errorf("tag:%s不存在", t)
				tx.Rollback() // 回滚事务
				return "", code.ErrorTagNotExit
			}
			articleTagModels = append(articleTagModels, model.ArticleTagModel{
				ArticleID:    id,
				ArticleTitle: articleTitle,
				TagTitle:     t,
			})
		}
		err = tx.Table("article_tag_models").CreateInBatches(articleTagModels, 100).Error
		if err != nil {
			global.Log.Errorf("article_tag_models update err: %s\n", err.Error())
			tx.Rollback() // 回滚事务
			return "", err
		}
	}
	err = tx.Commit().Error
	if err != nil {
		global.Log.Errorf("tx.Commit().Error err: %s\n", err.Error())
		tx.Rollback() // 回滚事务
		return "", err
	}
	return code.StrUpdateSucceed, nil
}
