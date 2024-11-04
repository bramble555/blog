package logic

import (
	"math/rand"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/bramble555/blog/dao/es"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/russross/blackfriday"
)

func UploadArticles(claims *pkg.MyClaims, pa *model.ParamArticle) (string, error) {
	userID := claims.ID
	username := claims.Username
	// 查询用户头像
	ud, err := user.GetUserDetail(userID)
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
	now := time.Now().Format("2006-01-02 15:04:05")
	article := model.ArticleModel{
		CreateTime: now,
		UpdateTime: now,
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
	return es.UploadArticles(article)
}
func GetArticlesList(pl *model.ParamList) (*[]model.ResponseArticle, error) {
	return es.GetArticlesList(pl)
}
