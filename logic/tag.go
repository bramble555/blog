package logic

import (
	"github.com/bramble555/blog/dao/mysql/tag"
	"github.com/bramble555/blog/model"
)

func CreateTags(tm *model.TagModel) (string, error) {
	return tag.CreateTags(tm)
}
func GetTags(pl *model.ParamList) ([]model.TagModel, error) {
	return tag.GetTags(pl)
}
func DeleteTagsList(pdl *model.ParamDeleteList) (string, error) {
	return tag.DeleteTagsList(pdl)
}
