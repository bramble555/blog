package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

func CreateTagsHandle(c *gin.Context) {
	var tm model.TagModel
	err := c.ShouldBindJSON(&tm)
	if err != nil {
		global.Log.Errorf("controller CreateTagsHandle ShouldBindJSON err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return
	}
	data, err := logic.CreateTags(&tm)
	if err != nil {
		global.Log.Errorf("controller CreateTagsHandle logic.CreateTags err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

func GetTagsListHandler(c *gin.Context) {
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetTagsListHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetTagsList(pl)
	if err != nil {
		global.Log.Errorf("controller GetTagsListHandler logic.GetTagsList err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSucceed(c, data)
}

func DeleteTagsListHandler(c *gin.Context) {
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteTagsListHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.DeleteTagsList(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteTagsListHandler logic.DeleteTagsList err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
