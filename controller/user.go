package controller

import (
	"errors"
	"time"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

// EmailLoginHandler 用户 用 email 或者用户名登录
func EmailLoginHandler(c *gin.Context) {
	var peu model.ParamEmailUser
	err := c.ShouldBindJSON(&peu)
	if err != nil {
		global.Log.Errorf("controller EmailLoginHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务处理
	token, err := logic.UsernameLogin(&peu)
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
func GetUserListHandler(c *gin.Context) {
	// 根据 token，获取 用户权限
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	role := claims.Role
	// ParamList 默认值
	var pl = model.ParamList{
		Page:  1,
		Size:  10,
		Order: model.OrderByTime,
	}
	err := c.ShouldBind(&pl)
	if err != nil {
		global.Log.Errorf("controller GetUserListHandler err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetUserList(role, &pl)
	if err != nil {
		ResponseError(c, CodeServerBusy)
	}
	ResponseSucceed(c, data)
}

// UpdateUserRoleHandler 管理员更新用户权限(其余均不更新)
func UpdateUserRoleHandler(c *gin.Context) {
	puur := model.ParamUpdateUserRole{}
	err := c.ShouldBindJSON(&puur)
	if err != nil {
		global.Log.Errorf("controller UpdateUserRoleHandler err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.UpdateUserRole(&puur)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}

// UpdateUserPwdHandler 用户更新用户新密码,需要旧密码
func UpdateUserPwdHandler(c *gin.Context) {
	puup := model.ParamUpdateUserPwd{}
	err := c.ShouldBindJSON(&puup)
	if err != nil {
		global.Log.Errorf("controller UpdateUserPwdHandler err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	var data string
	data, err = logic.UpdateUserPwd(&puup, claims.ID)
	if err != nil {
		ResponseErrorWithData(c, CodeServerBusy, data)
		return
	}
	ResponseSucceed(c, data)
}

// LogoutHandler 防止用户在注销后使用旧的 token 进行操作
// 注销是用户主动。(退出是用户被动，token 过期了)
func LogoutHandler(c *gin.Context) {
	// 获取 token
	token := c.Request.Header.Get("token")
	// 获取 MyClaims
	_cliams, _ := c.Get("claims")
	claims := _cliams.(*pkg.MyClaims)
	exp := time.Unix(claims.ExpiresAt, 0)
	now := time.Now()
	diff := exp.Sub(now)
	err := logic.Logout(token, diff)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, "注销成功")
}
