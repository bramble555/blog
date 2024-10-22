package controller

import (
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
