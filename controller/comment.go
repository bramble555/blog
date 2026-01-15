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
	uSN := claims.SN
	pc := model.ParamPostComment{}
	// 父评论默认为 -1
	pc.ParentCommentSN = -1
	err := c.ShouldBindJSON(&pc)
	if err != nil {
		global.Log.Errorf("controller PostArticleCommentsHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.PostArticleComments(uSN, &pc)
	if err != nil {
		global.Log.Errorf("controller PostArticleCommentsHandler logic.PostArticleComments err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func GetArticleCommentsHandler(c *gin.Context) {
	pcl := model.ParamCommentList{}
	err := c.ShouldBindQuery(&pcl)
	if err != nil {
		global.Log.Errorf("controller GetArticleCommentsHandler ShouldBindQuery err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var uSN int64
	_claims, exists := c.Get("claims")
	if exists {
		claims := _claims.(*pkg.MyClaims)
		uSN = claims.SN
	}
	data, err := logic.GetArticleComments(&pcl, uSN)
	if err != nil {
		global.Log.Errorf("controller GetArticleCommentsHandler logic.GetArticleComments err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
func DeleteArticleCommentsHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uSN := claims.SN
	role := claims.Role
	psn := model.ParamSN{}
	err := c.ShouldBindJSON(&psn)
	if err != nil {
		global.Log.Errorf("controller DeleteArticleCommentsHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	// We need ArticleSN to update the count, but frontend only sends sn.
	// We will look it up in the logic layer.
	data, err := logic.DeleteArticleComments(uSN, role, &psn, 0)
	if err != nil {
		global.Log.Errorf("controller DeleteArticleCommentsHandler logic.DeleteArticleComments err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
