package logic

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/bramble555/blog/dao/mysql/article"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/digg"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"github.com/bramble555/blog/pkg"
)

func UploadArticles(claims *pkg.MyClaims, pa *model.ParamArticle, bannerList *[]model.BannerModel) (string, error) {
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
	rawContent := pa.Content
	safeHTML := pkg.MarkdownToHTML(rawContent)

	if pa.Abstract == "" {
		c := []rune(rawContent)
		if len(c) >= 100 {
			pa.Abstract = string(c[:100])
		} else {
			pa.Abstract = string(c)
		}
	}
	// 随机选择一张图片作为为张图片
	rb := bannerList
	n := len(*rb)
	var bannerUrl string
	if n > 0 {
		// 设置随机种子（均匀分布，公平且高效）
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		idx := r.Intn(n)
		selected := (*rb)[idx]
		pa.BannerSN = selected.SN
		bannerUrl = "/uploads/file/" + selected.Name
		global.Log.Infof("selected banner sn: %d", pa.BannerSN)
	} else {
		// 无可用 banner 时，降级为无封面，避免异常
		pa.BannerSN = 0
		// 默认图片
		bannerUrl = "https://i.postimg.cc/nzwnF92j/dog.jpg"
		global.Log.Warnf("no banner available, article will use default cover")
	}

	tagList := pkg.ParseTagsStringSlice(pa.Tags)

	am := model.ArticleModel{
		Title:      pa.Title,
		Abstract:   pa.Abstract,
		Content:    rawContent,
		Tags:       ctype.ArrayString(tagList),
		BannerSN:   pa.BannerSN,
		BannerUrl:  bannerUrl,
		UserSN:     claims.SN,
		Username:   username,
		UserAvatar: ud.Avatar,
	}
	_, err = article.UploadArticles(&am)
	if err != nil {
		return "", err
	}
	response := map[string]any{
		"raw_content":    rawContent,
		"parsed_content": safeHTML,
		"status":         "ok",
	}
	//
	b, _ := json.Marshal(response)
	return string(b), nil
}

func GetArticlesListByParam(paq *model.ParamArticleQuery, uSN int64) (*model.ResponseArticleList, error) {
	return article.GetArticlesListByParam(paq, uSN)
}
func GetArticlesDetail(sn string, uSN int64) (*model.ArticleModel, error) {
	am, err := article.GetArticlesDetail(sn)
	if err != nil {
		return nil, err
	}
	// Parse Markdown
	am.ParsedContent = pkg.MarkdownToHTML(am.Content)

	// 检查用户是否登录,俩者页面是不一样的
	if uSN != 0 {
		// 检查用户是否收藏
		isCollect, err := article.IsUserCollect(uSN, am.SN)
		if err == nil && isCollect == true {
			am.IsCollect = isCollect
		}
		// 检查用户是否点赞
		isDigg, err := digg.IsUserDigg(uSN, am.SN)
		if err == nil && isDigg == true {
			am.IsDigg = isDigg
		}
	}
	return am, nil
}

func GetArticlesCalendar() (map[string]int, error) {
	return article.GetArticlesCalendar()
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
		return "", nil
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
