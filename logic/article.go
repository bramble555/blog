package logic

import (
	"math/rand"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/bramble555/blog/dao"
	"github.com/bramble555/blog/dao/es"
	"github.com/bramble555/blog/dao/mysql/article"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/russross/blackfriday"
)

func UploadArticles(claims *pkg.MyClaims, pa *model.ParamArticle) (string, error) {
	// 判断标题是否存在
	ok, err := article.TitleIsExist(pa.Title)
	if err != nil {
		return "", err
	}
	if ok {
		return "", code.ErrorTitleExit
	}
	userID := claims.ID
	username := claims.Username
	// 查询用户头像
	ud, err := user.GetUserDetailByID(userID)
	if err != nil {
		return "", err
	}
	// 处理文本 删除 script 标签
	ct := blackfriday.MarkdownCommon([]byte(pa.Content))
	// 查看是否有 script 标签
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(ct)))
	nodes := doc.Find("script").Nodes

	if len(nodes) > 0 {
		doc.Find("script").Remove()
		convert := md.NewConverter("", true, nil)
		html, _ := doc.Html()
		markdown, _ := convert.ConvertString(html)
		pa.Content = markdown

	}
	// 用户是否传入简介
	if pa.Abstract == "" {
		c := []rune(pa.Content)
		if len(c) >= 100 {
			pa.Abstract = string(c[:100])
		}
		pa.Abstract = string(c)
	}
	// 用户是否传入图片，如果没有，那就随机选择一张图片
	rb, err := GetBannerDetail()
	if err != nil {
		return "", err
	}
	n := len(*rb)
	IDList := make([]uint, n)
	for i, v := range *rb {
		IDList[i] = v.ID
	}
	// 设置随机种子
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	// 从 IDList 中随机选择一个 ID
	randomIndex := r.Intn(n) // 生成一个 [0, n) 范围内的随机数
	pa.BannerID = IDList[randomIndex]
	bannerUrl := (*rb)[randomIndex].Name
	// 组装数据
	am := model.ArticleModel{
		Title:      pa.Title,
		Abstract:   pa.Abstract,
		Content:    pa.Content,
		Category:   pa.Category,
		Source:     pa.Source,
		Link:       pa.Link,
		Tags:       pa.Tags,
		BannerID:   pa.BannerID,
		BannerUrl:  bannerUrl,
		UserID:     userID,
		Username:   username,
		UserAvatar: ud.Avatar,
	}
	return article.UploadArticles(&am)
}



func GetArticlesListByParam() dao.ArticleQueryService {
	if global.Config.ES.Enable {
		return &es.ESArticleQueryService{} // 返回 Elasticsearch 查询服务
	}
	return &article.MySQLArticleQueryService{} // 返回 MySQL 查询服务
}
func GetArticlesDetail(id string) (*model.ArticleModel, error) {
	return article.GetArticlesDetail(id)
}

func GetArticlesCalendar() (*map[string]int, error) {
	return article.GetArticlesCalendar()
}
func GetArticlesTagsList(paq *model.ParamList) (*[]model.ResponseArticleTags, error) {
	return article.GetArticlesTagsList(paq)
}
func UpdateArticles(id uint, uf map[string]any) (string, error) {
	ok, err := article.IDExist(id)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorIDNotExit
	}
	return article.UpdateArticles(id, uf)
}
func DeleteArticlesList(pdl *model.ParamDeleteList) (string, error) {
	// 检查 IDList 是否为空
	if len(pdl.IDList) == 0 {
		return "", code.ErrorIDNotExit
	}
	// 查询 IDList 是否存在
	ok, err := article.IDListExist(pdl)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorIDNotExit
	}
	return article.DeleteArticlesList(pdl)
}
func PostArticleCollect(uID uint, articleID uint) (string, error) {
	// 查询 articleID 是否存在
	ok, err := article.IDExist(articleID)
	if err != nil {
		return "", err
	}
	if !ok {
		global.Log.Errorf("文章 ID:%d不存在\n", articleID)
		return "", code.ErrorIDExit
	}
	return article.PostArticleCollect(uID, articleID)
}
func GetArticleCollect(uID uint) ([]model.ResponseArticle, error) {
	return article.GetArticleCollect(uID)
}
func DeleteArticleCollect(uID uint, pdl *model.ParamDeleteList) (string, error) {
	// 检查用户已经收藏的文章
	count, err := article.GetUserCollectsCount(uID, pdl.IDList)
	if err != nil {
		return "", nil
	}
	if int(count) != len(pdl.IDList) {
		global.Log.Errorf("IDList 不存在")
		return "", code.ErrorIDExit
	}
	return article.DeleteArticleCollect(uID, pdl.IDList)
}
