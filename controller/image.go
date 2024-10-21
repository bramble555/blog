package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model/image"
	"github.com/gin-gonic/gin"
)

func ImageHandler(c *gin.Context) {
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
	data, err = logic.UploadImage(c, fileList)
	if err != nil {
		global.Log.Errorf("Logic UploadImage  [images] err:%s\n", err.Error())
	}
	ResponseSucceed(c, data)
}
