package controller

import (
	"strings"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

func CreateAdvertHandle(c *gin.Context) {
	var ad model.AdvertModel
	err := c.ShouldBindJSON(&ad)
	if err != nil {
		global.Log.Errorf("controller CreateAdvertHandle ShouldBindQuery err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}
	data, err := logic.CreateAdvert(&ad)
	if err != nil {
		global.Log.Errorf("logic CreateAdvertHandle ShouldBindQuery err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
func GetAdvertListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetAdvertListHandler err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	// 判断 referer 是否包含 admin，如果是，返回，如果不是，就不返回了
	referer := c.GetHeader("referer")
	isShow := false
	if strings.Contains(referer, "admin") {
		isShow = true
	}
	data, err := logic.GetAdvertList(pl, isShow)
	if err != nil {
		global.Log.Errorf("Logic GetImageList err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSucceed(c, data)
}
func DeleteAdvertListHander(c *gin.Context) {
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("DeleteHanderListHander ShouldBindQuery err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}
	var data string
	data, err = logic.DeleteAdvertList(&pdl)
	if err != nil {
		global.Log.Errorf("logic DeleteAdvertList err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
