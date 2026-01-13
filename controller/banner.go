package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

func UploadBannersHandler(c *gin.Context) {
	// 获取图片并且对参数进行处理
	form, err := c.MultipartForm()
	if err != nil {
		global.Log.Errorf("controller UploadBannersHandler c.MultipartForm err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}
	fileList, ok := form.File["images"]
	if !ok {
		global.Log.Errorf("controller UploadBannersHandler form.File[images] not found\n")
		ResponseError(c, CodeInvalidParam)
		return
	}
	data := new([]model.FileUploadResponse)
	data, err = logic.UploadImages(c, fileList)
	if err != nil {
		global.Log.Errorf("controller UploadBannersHandler logic.UploadImages err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetBannerListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetBannerListHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetBannerList(pl)
	if err != nil {
		global.Log.Errorf("controller GetBannerListHandler logic.GetBannerList err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSucceed(c, data)
}
func GetBannerDetailHandler(c *gin.Context) {
	data, err := logic.GetBannerDetail()
	if err != nil {
		global.Log.Errorf("controller GetBannerDetailHandler logic.GetBannerDetail err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func DeleteBannerListHander(c *gin.Context) {
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteBannerListHander ShouldBindJSON err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}
	var data string
	data, err = logic.DeleteBannerList(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteBannerListHander logic.DeleteBannerList err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
