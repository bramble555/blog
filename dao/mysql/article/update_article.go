package article

import (
	"strconv"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"github.com/bramble555/blog/pkg/convert"
)

// UpdateArticles 更新文章信息
func UpdateArticles(sn int64, uf map[string]any) (string, error) {
	// 开始数据库事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Log.Errorf("article_models start transaction failed: %v", tx.Error)
		return "", tx.Error
	}

	// 如果存在“content”字段，则提取其摘要并保存在“abstract”字段
	if content, ok := uf["content"].(string); ok {
		runes := []rune(content)
		if len(runes) > 100 {
			uf["abstract"] = string(runes[:100]) // 如果内容长度大于100个字符，只取前100个字符作为摘要
		} else {
			uf["abstract"] = content // 如果内容小于或等于100个字符，摘要即为内容
		}
	}

	// 处理标签字段
	var uniqueTags []string
	var hasTags bool
	if tagsInterface, exists := uf["tags"]; exists {
		switch v := tagsInterface.(type) {
		case string:
			hasTags = true
			uniqueTags = convert.ParseTagsStringSlice(v) // 将标签字符串解析为切片
		}
		if hasTags {
			uf["tags"] = ctype.ArrayString(uniqueTags) // 将标签赋值为一个字符串数组
		}
	}

	// 处理 banner_sn 字段
	if bStr, ok := uf["banner_sn"].(string); ok {
		if bSN, err := strconv.ParseInt(bStr, 10, 64); err == nil {
			uf["banner_sn"] = bSN // 如果是有效的字符串数字，转换为 int64 并更新
		}
	}

	// 更新文章信息
	if err := tx.Model(&model.ArticleModel{}).Where("sn = ?", sn).Updates(uf).Error; err != nil {
		global.Log.Errorf("article_models update failed: %v", err)
		tx.Rollback() // 如果更新失败，回滚事务
		return "", err
	}

	// 检查是否有更新标题和标签
	articleTitleVal, hasTitle := uf["title"]
	if !hasTitle && !hasTags {
		if err := tx.Commit().Error; err != nil {
			return "", err
		}
		return code.StrUpdateSucceed, nil
	}

	// 如果更新了标题，则获取标题值
	var articleTitle string
	if hasTitle {
		articleTitle, _ = articleTitleVal.(string)
	} else {
		// 如果没有标题，则从数据库中读取文章原始标题
		if err := tx.Model(&model.ArticleModel{}).Select("title").Where("sn = ?", sn).Scan(&articleTitle).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	// 处理标签更新
	if hasTags {
		// 删除现有的标签关联
		if err := tx.Where("article_sn = ?", sn).Delete(&model.ArticleTagModel{}).Error; err != nil {
			global.Log.Errorf("delete old tags failed: %v", err)
			tx.Rollback()
			return "", err
		}

		// 检查标签是否存在于数据库中
		if len(uniqueTags) > 0 {
			var dbTags []string
			if err := tx.Model(&model.TagModel{}).Where("title IN ?", uniqueTags).Pluck("title", &dbTags).Error; err != nil {
				tx.Rollback()
				return "", err
			}

			// 如果数据库中没有某些现在需要上传标签，返回 code.ErrorTagNotExist
			if len(dbTags) != len(uniqueTags) {
				dbTagMap := make(map[string]struct{}, len(dbTags))
				for _, t := range dbTags {
					dbTagMap[t] = struct{}{}
				}
				for _, reqTag := range uniqueTags {
					if _, ok := dbTagMap[reqTag]; !ok {
						global.Log.Errorf("tag not found in tag_models: %s", reqTag)
						tx.Rollback()
						return "", code.ErrorTagNotExist
					}
				}
			}

			// 批量插入新的标签关系
			newRelations := make([]model.ArticleTagModel, 0, len(uniqueTags))
			for _, tag := range uniqueTags {
				newRelations = append(newRelations, model.ArticleTagModel{
					ArticleSN:    sn,
					TagTitle:     tag,
					ArticleTitle: articleTitle,
				})
			}
			// 分批次插入,如果一个文章有很多标签,一批次插入5个
			if err := tx.CreateInBatches(newRelations, 5).Error; err != nil {
				global.Log.Errorf("batch insert tags failed: %v", err)
				tx.Rollback()
				return "", err
			}
		}
	} else if hasTitle {
		// 如果只有标题更新，更新文章标签模型中的标题字段
		if err := tx.Model(&model.ArticleTagModel{}).Where("article_sn = ?", sn).
			Update("article_title", articleTitle).Error; err != nil {
			global.Log.Errorf("sync title to tags failed: %v", err)
			tx.Rollback()
			return "", err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return code.StrUpdateSucceed, nil
}
