package controller

import (
	"errors"
	"strconv"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

// UploadArticleHandler 上传文章
func UploadArticleHandler(c *gin.Context) {
	var pa model.ParamArticle
	err := c.ShouldBindJSON(&pa)
	if err != nil {
		global.Log.Errorf("controller UploadArticlesHandler err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)

	// 获取 banner 列表，用于随机选择封面
	pl := &model.ParamList{Page: 1, Size: 100}
	bannerList, err := logic.GetBannerList(pl)
	if err != nil {
		global.Log.Errorf("controller UploadArticlesHandler logic.GetBannerList err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}

	data, err := logic.UploadArticles(claims, &pa, bannerList)
	if err != nil {
		if errors.Is(err, code.ErrorTitleExit) {
			ResponseError(c, CodeTitleExist)
			return
		}
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

// GetArticlesListHandler 获取文章列表，如果有 title ，page 等字段也会根据其进行搜索
func GetArticlesListHandler(c *gin.Context) {
	paq := new(model.ParamArticleQuery)
	err := c.ShouldBindQuery(paq)
	if err != nil {
		global.Log.Errorf("controller GetArticlesListHandler ShouldBindQuery err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	if paq.Order == "" {
		paq.Order = model.OrderByTime
	}
	if paq.Page <= 0 {
		paq.Page = 1
	}
	if paq.Size <= 0 {
		paq.Size = 10
	}

	var uSN int64
	_claims, exists := c.Get("claims")
	if exists {
		claims := _claims.(*pkg.MyClaims)
		uSN = claims.SN
	}

	data, err := logic.GetArticlesListByParam(paq, uSN)
	if err != nil {
		global.Log.Errorf("controller GetArticlesListHandler logic.GetArticlesListByParam err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetArticlesDetailHandler(c *gin.Context) {
	sn := c.Param("sn")
	var uSN int64
	_claims, exists := c.Get("claims")
	if exists {
		claims := _claims.(*pkg.MyClaims)
		uSN = claims.SN
	}
	data, err := logic.GetArticlesDetail(sn, uSN)
	if err != nil {
		global.Log.Errorf("controller GetArticlesDetailHandler logic.GetArticlesDetail err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	ResponseSucceed(c, data)
}

func GetArticlesCalendarHandler(c *gin.Context) {
	data, err := logic.GetArticlesCalendar()
	if err != nil {
		global.Log.Errorf("controller GetArticlesCalendarHandler logic.GetArticlesCalendar err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetArticlesTagsListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetArticlesTagsListHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetArticlesTagsList(pl)
	if err != nil {
		global.Log.Errorf("controller GetArticlesTagsListHandler logic.GetArticlesTagsList err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func UpdateArticlesHandler(c *gin.Context) {
	// http://localhost:8080/articles/1
	articleSN := c.Param("sn")
	// 将 articleSN 转换为整型
	sn, err := strconv.ParseInt(articleSN, 10, 64)
	if err != nil {
		global.Log.Errorf("controller UpdateArticlesHandler strconv.ParseInt err:%s\n", err.Error())
		ResponseError(c, CodeInvalidID)
		return
	}
	uf := make(map[string]any)
	if err := c.ShouldBindJSON(&uf); err != nil {
		global.Log.Errorf("controller UpdateArticlesHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.UpdateArticles(int64(sn), uf)
	if err != nil {
		global.Log.Errorf("controller UpdateArticlesHandler logic.UpdateArticles err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
func DeleteArticlesListHandler(c *gin.Context) {
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteArticlesListHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.DeleteArticlesList(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteArticlesListHandler logic.DeleteArticlesList err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
func PostArticleCollectHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uSN := claims.SN
	psn := model.ParamSN{}
	err := c.ShouldBindJSON(&psn)
	if err != nil {
		global.Log.Errorf("controller PostArticleCollectHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.PostArticleCollect(uSN, psn.SN)
	if err != nil {
		global.Log.Errorf("controller PostArticleCollectHandler logic.PostArticleCollect err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetArticleCollectHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uSN := claims.SN
	data, err := logic.GetArticleCollect(uSN)
	if err != nil {
		global.Log.Errorf("controller GetArticleCollectHandler logic.GetArticleCollect err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func DeleteArticleCollectHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uSN := claims.SN
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteArticleCollectHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.DeleteArticleCollect(uSN, &pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteArticleCollectHandler logic.DeleteArticleCollect err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
