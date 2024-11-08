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
func GetArticlesListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetArticlesList(pl)
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
	uf := model.UpdatedFields{}
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
