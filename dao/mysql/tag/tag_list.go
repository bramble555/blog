package tag

import (
	"fmt"

	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

func GetTagsList(pl *model.ParamList) ([]model.TagModel, error) {
	// 调用 GetTableList 时，只有在 condition 不为空时才传递条件和参数
	tml, err := mysql.GetTableList[model.TagModel]("tag_models", pl, "")
	if err != nil {
		return nil, err
	}
	return tml, nil
}

// DeleteTagsList 删除 tags 列表，并同步删除文章与标签的关联关系
func DeleteTagsList(pdl *model.ParamDeleteList) (string, error) {
	// 转换 SNList 为 []int64
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("DeleteTagsList StringSliceToInt64Slice err: %s\n", err.Error())
		return "", err
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 获取要删除的标签信息，主要是为了获取 title
	var tags []model.TagModel
	if err := tx.Table("tag_models").Where("sn IN ?", snList).Find(&tags).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("DeleteTagsList find tags err: %s\n", err.Error())
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

	// 2. 删除文章与标签的关联关系
	if err := tx.Table("article_tag_models").Where("tag_title IN ?", tagTitles).Delete(nil).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("DeleteTagsList delete article_tag_models err: %s\n", err.Error())
		return "", err
	}

	// 3. 删除标签本身
	result := tx.Table("tag_models").Where("sn IN ?", snList).Delete(&tags)
	if result.Error != nil {
		tx.Rollback()
		global.Log.Errorf("DeleteTagsList delete tag_models err: %s\n", result.Error)
		return "", result.Error
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("DeleteTagsList commit err: %s\n", err.Error())
		return "", err
	}

	return fmt.Sprintf("成功删除 %d 个标签及其关联关系", result.RowsAffected), nil
}
