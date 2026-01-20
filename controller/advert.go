package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

// UploadAdvertImagesHandler 处理广告图片上传
func UploadAdvertImagesHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		global.Log.Errorf("controller UploadAdvertImagesHandler c.MultipartForm err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}

	fileList, ok := form.File["images"]
	if !ok {
		global.Log.Errorf("controller UploadAdvertImagesHandler form.File[images] not found\n")
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.UploadAdvertImages(c, fileList)
	if err != nil {
		global.Log.Errorf("controller UploadAdvertImagesHandler logic.UploadAdvertImages err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSucceed(c, data)
}

// GetAdvertListHandler 管理员可以获取广告列表
func GetAdvertListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetAdvertListHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	// is_show=true -> 前台只显示可见项；未传参 -> 管理后台显示全部
	qs := c.Query("is_show")
	// DAO 中 isShow=false 时才添加过滤条件 is_show=true
	// 因此当 qs=="true" 时传 false，其他情况传 true（返回全部）
	isShow := true
	if qs == "true" {
		isShow = false
	}
	data, err := logic.GetAdvertList(pl, isShow)
	if err != nil {
		global.Log.Errorf("controller GetAdvertListHandler logic.GetAdvertList err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSucceed(c, data)
}
func DeleteAdvertListHandler(c *gin.Context) {
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteAdvertListHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.DeleteAdvertList(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteAdvertListHandler logic.DeleteAdvertList err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

func UpdateAdvertShowHandler(c *gin.Context) {
	var p model.ParamUpdateAdvertShow
	err := c.ShouldBindJSON(&p)
	if err != nil {
		global.Log.Errorf("controller UpdateAdvertShowHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.UpdateAdvertShow(&p)
	if err != nil {
		global.Log.Errorf("controller UpdateAdvertShowHandler logic.UpdateAdvertShow err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
