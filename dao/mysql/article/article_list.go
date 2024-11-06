package article

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func GetArticlesList(pl *model.ParamList) (*[]model.ResponseArticle, error) {
	offset := (pl.Page - 1) * pl.Size
	res := []model.ResponseArticle{}
	err := global.DB.Table("article_models").
		Select("id, create_time, update_time, title, abstract, look_count, comment_count, digg_count, collects_count, category, tags, banner_id, banner_url, user_id, username, user_avatar").
		Order(pl.Order).
		Limit(pl.Size).
		Offset(offset).
		Find(&res).Error
	if err != nil {
		global.Log.Errorf("select err:%s\n", err.Error())
		return nil, err
	}
	return &res, nil
}
func GetArticlesDetail(id string) (*model.ArticleModel, error) {
	am := model.ArticleModel{}
	err := global.DB.Table("article_models").
		Where("id = ?", id).
		First(&am).Error

	if err != nil {
		// 错误处理，输出日志
		global.Log.Errorf("select err:%s\n", err.Error())
		return nil, err
	}
	return &am, nil
}
