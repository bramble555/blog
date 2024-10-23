package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/image"
	"github.com/gin-gonic/gin"
)

func UploadBannerHandler(c *gin.Context) {
	// 获取图片并且对参数进行处理
	form, err := c.MultipartForm()
	if err != nil {
		global.Log.Errorf("Controller ImageHandler err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
	}
	fileList, ok := form.File["images"]
	if !ok {
		global.Log.Errorf("Controller ImageHandler formFile[images] err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
	}
	data := new([]image.FileUploadResponse)
	data, err = logic.UploadImages(c, fileList)
	if err != nil {
		global.Log.Errorf("Logic UploadImage  [images] err:%s\n", err.Error())
	}
	ResponseSucceed(c, data)
}
func GetBannerListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetBannerListHandler err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	data, err := logic.GetBannerList(pl)
	if err != nil {
		global.Log.Errorf("Logic GetBannerListHandler err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSucceed(c, data)
}
func DeleteBannerListHander(c *gin.Context) {
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("DeleteHanderListHander ShouldBindQuery err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}
	var data string
	data, err = logic.DeleteBannerList(&pdl)
	if err != nil {
		global.Log.Errorf("logic DeleteBannerList err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
