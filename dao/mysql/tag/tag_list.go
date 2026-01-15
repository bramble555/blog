package tag

import (
	"fmt"
	"strings"

	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

// GetTagsList 获取 tags 列表
func GetTagsList(pl *model.ParamList) ([]model.TagModel, error) {
	tml, err := mysql.GetTableList[model.TagModel]("tag_models", pl, "")
	if err != nil {
		return nil, err
	}
	return tml, nil
}

// DeleteTagsList 删除 tags 列表，并同步删除文章与标签的关联关系
func DeleteTagsList(pdl *model.ParamDeleteList) (string, error) {
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("DeleteTagsList err: %s", err.Error())
		return "", err
	}

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 获取要删除的标签信息
	var tags []model.TagModel
	if err := tx.Where("sn IN ?", snList).Find(&tags).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	if len(tags) == 0 {
		tx.Rollback()
		return "没有找到要删除的标签", nil
	}

	var tagTitles []string
	for _, tag := range tags {
		tagTitles = append(tagTitles, tag.Title)
	}

	// 2. 在删除前，找出哪些文章受到了影响
	var affectedArticleSNs []int64
	tx.Table("article_tag_models").
		Where("tag_title IN ?", tagTitles).
		Distinct("article_sn").
		Pluck("article_sn", &affectedArticleSNs)

	// 3. 删除关联关系
	if err := tx.Table("article_tag_models").Where("tag_title IN ?", tagTitles).Delete(nil).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// 4. 删除标签本身
	if err := tx.Where("sn IN ?", snList).Delete(&model.TagModel{}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// 5. 更新文章表中的冗余 tags 字符串
	if len(affectedArticleSNs) > 0 {
		for _, aSN := range affectedArticleSNs {
			// 重新查一下这篇文章现在还剩下哪些标签
			var currentTitles []string
			tx.Table("article_tag_models").
				Where("article_sn = ?", aSN).
				Pluck("tag_title", &currentTitles)

			// 将剩下的标签数组转为逗号分隔的字符串
			newTagsStr := strings.Join(currentTitles, ",")
			if newTagsStr == "" {
				newTagsStr = "无标签待更新"
			}
			// 更新到文章表
			if err := tx.Table("article_models").
				Where("sn = ?", aSN).
				Update("tags", newTagsStr).Error; err != nil {
				tx.Rollback()
				global.Log.Errorf("Update article tags redundancy err: %s", err.Error())
				return "", err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("成功删除 %d 个标签，并同步更新了 %d 篇文章的标签显示", len(tags), len(affectedArticleSNs)), nil
}
