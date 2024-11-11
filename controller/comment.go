package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

func PostArticleCommentsHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uID := claims.ID
	pc := model.ParamPostComment{}
	// 父评论默认为 -1
	pc.ParentCommentID = -1
	err := c.ShouldBindJSON(&pc)
	if err != nil {
		global.Log.Errorf("ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.PostArticleComments(uID, &pc)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetArticleCommentsHandler(c *gin.Context) {
	pcl := model.ParamCommentList{}
	err := c.ShouldBindQuery(&pcl)
	if err != nil {
		global.Log.Errorf("ShouldBindQuery err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetArticleComments(&pcl)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
