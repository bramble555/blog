package controller

import (
	"errors"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

// EmailLoginHandler 用户 email 或者用户名登录
func EmailLoginHandler(c *gin.Context) {
	var peu model.ParamEmailUser
	err := c.ShouldBindJSON(&peu)
	if err != nil {
		global.Log.Errorf("controller EmailLoginHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务处理
	token, err := logic.EmailLogin(&peu)
	if err != nil {
		global.Log.Errorf("Login with invaild params:%s\n", err.Error())
		if err == errors.New("用户名不存在") {
			ResponseError(c, CodeUserExist)
			return
		} else if err == errors.New("密码错误") {
			ResponseError(c, CodeInvalidPassword)
			return
		} else {
			ResponseErrorWithData(c, CodeInvalidParam, err.Error())
			return
		}
	}
	ResponseSucceed(c, token)
}
