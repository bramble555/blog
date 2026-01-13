package logic

import (
	"math/rand"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/bramble555/blog/dao"
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
	username := claims.Username
	// 查询用户头像
	ud, err := user.GetUserDetailBySN(claims.SN)
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
	SNList := make([]int64, n)
	for i, v := range *rb {
		SNList[i] = v.SN
	}
	// 设置随机种子
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	// 从 SNList 中随机选择一个 SN
	randomIndex := r.Intn(n) // 生成一个 [0, n) 范围内的随机数
	pa.BannerSN = SNList[randomIndex]
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
		BannerSN:   pa.BannerSN,
		BannerUrl:  bannerUrl,
		UserSN:     int64(claims.SN), // Cast int64 to int64 for database field
		Username:   username,
		UserAvatar: ud.Avatar,
	}
	return article.UploadArticles(&am)
}

func GetArticlesListByParam() dao.ArticleQueryService {
	return &article.MySQLArticleQueryService{} // 返回 MySQL 查询服务
}
func GetArticlesDetail(sn string, uSN int64) (*model.ArticleModel, error) {
	am, err := article.GetArticlesDetail(sn)
	if err != nil {
		return nil, err
	}
	if uSN != 0 {
		// 检查用户是否收藏
		isCollect, err := article.IsUserCollect(uSN, am.SN)
		if err == nil {
			am.IsCollect = isCollect
		}
	}
	return am, nil
}

func GetArticlesCalendar() (*map[string]int, error) {
	return article.GetArticlesCalendar()
}
func GetArticlesTagsList(paq *model.ParamList) (*[]model.ResponseArticleTags, error) {
	return article.GetArticlesTagsList(paq)
}
func UpdateArticles(sn int64, uf map[string]any) (string, error) {
	ok, err := article.CheckSNExist(sn)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExit
	}
	return article.UpdateArticles(sn, uf)
}
func DeleteArticlesList(pdl *model.ParamDeleteList) (string, error) {
	// 检查 SNList 是否为空
	if len(pdl.SNList) == 0 {
		return "", code.ErrorSNNotExit
	}
	// 查询 SNList 是否存在
	ok, err := article.SNListExist(pdl)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExit
	}
	return article.DeleteArticlesList(pdl)
}
func PostArticleCollect(uSN int64, articleSN int64) (string, error) {
	// 查询 articleSN 是否存在
	ok, err := article.CheckSNExist(articleSN)
	if err != nil {
		return "", err
	}
	if !ok {
		global.Log.Errorf("文章 SN:%d不存在\n", articleSN)
		return "", code.ErrorSNNotExit
	}
	return article.PostArticleCollect(uSN, articleSN)
}
func GetArticleCollect(uSN int64) ([]model.ResponseArticle, error) {
	return article.GetArticleCollect(uSN)
}
func DeleteArticleCollect(uSN int64, pdl *model.ParamDeleteList) (string, error) {
	// 转换 SNList 为 []int64
	snList, err := pkg.StringSliceToInt64Slice(pdl.SNList)
	if err != nil {
		global.Log.Errorf("DeleteArticleCollect StringSliceToInt64Slice err: %s\n", err.Error())
		return "", err
	}

	// 检查用户已经收藏的文章
	count, err := article.GetUserCollectsCount(uSN, snList)
	if err != nil {
		return "", err
	}
	if int(count) != len(snList) {
		global.Log.Errorf("SNList 不存在")
		return "", code.ErrorSNNotExit
	}
	return article.DeleteArticleCollect(uSN, snList)
}
