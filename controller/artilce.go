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

func UploadArticlesHandler(c *gin.Context) {
	var pa model.ParamArticle
	err := c.ShouldBindJSON(&pa)
	if err != nil {
		global.Log.Errorf("controller UploadArticlesHandler err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)

	data, err := logic.UploadArticles(claims, &pa)
	if err != nil {
		if errors.Is(err, code.ErrorTitleExit) {
			ResponseError(c, CodeTitleExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}

// GetArticlesListHandler 获取文章列表，如果有 title ，page 等字段也会根据其进行搜索
func GetArticlesListHandler(c *gin.Context) {
	paq := new(model.ParamArticleQuery)
	err := c.ShouldBindQuery(paq)
	if err != nil {
		global.Log.Errorf("ShouldBindQuery err:%s\n", err.Error())
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
	data, err := logic.GetArticlesList(paq)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetArticlesDetailHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := logic.GetArticlesDetail(id)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	ResponseSucceed(c, data)
}

func GetArticlesCalendarHandler(c *gin.Context) {
	data, err := logic.GetArticlesCalendar()
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetArticlesTagsListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetArticlesTagsList(pl)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func UpdateArticlesHandler(c *gin.Context) {
	// http://localhost:8080/articles/1
	articleID := c.Param("id")
	// 将 articleID 转换为整型
	id, err := strconv.Atoi(articleID)
	if err != nil {
		global.Log.Errorf("id 有误 err:%s\n", err.Error())
		ResponseError(c, CodeInvalidID)
		return
	}
	uf := make(map[string]any)
	if err := c.ShouldBindJSON(&uf); err != nil {
		global.Log.Errorf("ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.UpdateArticles(uint(id), uf)
	if err != nil {
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
func DeleteArticlesListHandler(c *gin.Context) {
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("DeleteArticlesListHandler ShouldBindQuery err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.DeleteArticlesList(&pdl)
	if err != nil {
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
func PostArticleCollectHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uID := claims.ID
	pi := model.ParamID{}
	err := c.ShouldBindJSON(&pi)
	if err != nil {
		global.Log.Errorf("ShouldBindJSON err%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.PostArticleCollect(uID, pi.ID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetArticleCollectHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uID := claims.ID
	data, err := logic.GetArticleCollect(uID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func DeleteArticleCollectHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uID := claims.ID
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("DeleteArticleCollectHandler ShouldBindQuery err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.DeleteArticleCollect(uID, &pdl)
	if err != nil {
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
func PostArticleDigHandler(c *gin.Context) {
	pi := model.ParamID{}
	err := c.ShouldBindJSON(&pi)
	if err != nil {
		global.Log.Errorf("ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.PostArticleDig(pi.ID)
	if err != nil {
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
