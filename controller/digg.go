package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

// PostArticleDigHandler 实现了用户对文章点赞或取消点赞
func PostArticleDigHandler(c *gin.Context) {
	pi := model.ParamSN{}
	err := c.ShouldBindJSON(&pi)
	if err != nil {
		global.Log.Errorf("controller PostArticleDigHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uSN := claims.SN

	var data string
	data, err = logic.PostArticleDig(uSN, pi.SN)
	if err != nil {
		global.Log.Errorf("controller PostArticleDigHandler logic.PostArticleDig err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}

func PostArticleCommentsDiggHandler(c *gin.Context) {
	ps := model.ParamSN{}
	err := c.ShouldBindJSON(&ps)
	if err != nil {
		global.Log.Errorf("controller PostArticleCommentsDiggHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	uSN := claims.SN

	var data string
	data, err = logic.PostArticleCommentDig(uSN, ps.SN)
	if err != nil {
		global.Log.Errorf("controller PostArticleCommentsDiggHandler logic.PostArticleCommentDig err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSucceed(c, data)
}
