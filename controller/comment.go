package controller

import (
	"strconv"

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

// GetArticleCommentsHandler 获取文章评论列表
// 有文章 ID, 则是文章详情的函数, 没有文章 ID,是后台所有数据的函数
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

// DeleteCommentsHandler 删除评论
// 如果提供了 SNList, 则删除后台中的选中所有评论
// 如果提供了 SN, 则删除该评论下的所有子评论
func DeleteCommentsHandler(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uSN := claims.SN
	role := claims.Role
	var req model.ParamDeleteComment
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Errorf("controller DeleteArticleCommentsHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 如果提供了 SNList, 则删除后台中的选中所有评论
	if len(req.SNList) > 0 {
		data, err := logic.DeleteCommentsList(uSN, role, &model.ParamDeleteList{SNList: req.SNList})
		if err != nil {
			global.Log.Errorf("controller DeleteArticleCommentsHandler logic.DeleteCommentsList err:%s\n", err.Error())
			ResponseErrorWithData(c, CodeServerBusy, err.Error())
			return
		}
		ResponseSucceed(c, data)
		return
	}

	if req.SN == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	sn, err := strconv.ParseInt(req.SN, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.DeleteArticleComments(uSN, role, sn, 0)
	if err != nil {
		global.Log.Errorf("controller DeleteArticleCommentsHandler logic.DeleteArticleComments err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
