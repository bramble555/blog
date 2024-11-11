package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

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
func PostArticleCommentsDiggHandler(c *gin.Context) {
	pi := model.ParamID{}
	err := c.ShouldBindJSON(&pi)
	if err != nil {
		global.Log.Errorf("ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.PostArticleCommentDig(pi.ID)
	if err != nil {
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
