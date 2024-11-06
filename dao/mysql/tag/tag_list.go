package tag

import (
	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/model"
)

func GetAdvertList(pl *model.ParamList) ([]model.TagModel, error) {
	// 调用 GetTableList 时，只有在 condition 不为空时才传递条件和参数
	tml, err := mysql.GetTableList[model.TagModel]("tag_models", pl, "")
	if err != nil {
		return nil, err
	}
	return tml, nil
}

// DeleteTagsList 删除 tags 列表，但是这里没有实现同步删除文章
func DeleteTagsList(pdl *model.ParamDeleteList) (string, error) {
	return mysql.DeleteTableList[model.TagModel]("tag_models", pdl)
}
