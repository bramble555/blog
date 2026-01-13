package controller

import (
	"errors"
	"strconv"
	"time"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func PostBindEmailHandler(c *gin.Context) {
	var pbe model.ParamBindEmail
	// 获取 MyClaims
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)

	// 参数绑定
	err := c.ShouldBindJSON(&pbe)
	if err != nil {
		global.Log.Errorf("controller PostBindEmailHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 初始化 session
	session := sessions.Default(c)
	if pbe.Code == nil {
		// --- 发送验证码流程 ---
		code, err := pkg.SendEmail(pbe.Email)
		if err != nil {
			global.Log.Errorf("controller PostBindEmailHandler pkg.SendEmail err:%s", err.Error())
			ResponseError(c, CodeServerBusy)
			return
		}

		// 设置 session 值
		session.Set("valid_email", pbe.Email)
		session.Set("code", strconv.Itoa(code))
		// 保存 session
		session.Save()
		ResponseSucceed(c, "验证码发送成功")
		return
	}

	// --- 校验验证码流程 ---
	storedEmail := session.Get("valid_email")
	storedCode := session.Get("code")
	// 验证 email 是否一致
	if storedEmail != pbe.Email {
		global.Log.Errorf("controller PostBindEmailHandler email mismatch\n")
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 验证 code 是否正确
	if storedCode != *pbe.Code {
		global.Log.Errorf("controller PostBindEmailHandler code mismatch\n")
		ResponseError(c, CodeInvalidVerification)
		return
	}

	// 检查是否需要修改邮箱
	ud, err := user.GetUserDetailBySN(claims.SN)
	if err != nil {
		global.Log.Errorf("controller PostBindEmailHandler user.GetUserDetailBySN err:%s", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	// 检查是否与当前邮箱一致
	if ud.Email == pbe.Email {
		ResponseSucceed(c, "邮箱无需修改")
		return
	}

	// 更新邮箱（绑定或修改）
	err = logic.PostBindEmail(claims.SN, pbe.Email) // claims.SN is int64
	if err != nil {
		global.Log.Errorf("controller PostBindEmailHandler logic.PostBindEmail err:%s", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSucceed(c, "邮箱修改成功")
}

// UsernameLoginHandler 用户用用户名登录
func UsernameLoginHandler(c *gin.Context) {
	var peu model.ParamUsername
	err := c.ShouldBindJSON(&peu)
	if err != nil {
		global.Log.Errorf("controller UsernameLoginHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务处理
	data, err := logic.UsernameLogin(&peu)
	if err != nil {
		if errors.Is(err, code.ErrorUserNotExit) {
			ResponseErrorWithData(c, CodeUserNotExist, err)
			return
		} else if errors.Is(err, code.ErrorPasswordWrong) {
			ResponseErrorWithData(c, CodeInvalidPassword, err)
			return
		} else {
			ResponseError(c, CodeInvalidParam)
			return
		}
	}
	ResponseSucceed(c, data)
}
func GetUserListHandler(c *gin.Context) {
	// 根据 token，获取 用户权限
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	role := claims.Role
	// ParamList 默认值
	pl, err := validateListParams(c)
	if err != nil {
		global.Log.Errorf("controller GetUserListHandler validateListParams err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetUserList(role, pl)
	if err != nil {
		global.Log.Errorf("controller GetUserListHandler logic.GetUserList err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}

// UpdateUserRoleHandler 管理员更新用户权限(其余均不更新)
func UpdateUserRoleHandler(c *gin.Context) {
	puur := model.ParamUpdateUserRole{}
	err := c.ShouldBindJSON(&puur)
	if err != nil {
		global.Log.Errorf("controller UpdateUserRoleHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.UpdateUserRole(&puur)
	if err != nil {
		global.Log.Errorf("controller UpdateUserRoleHandler logic.UpdateUserRole err:%s\n", err.Error())
		if errors.Is(err, code.ErrorSNNotExit) {
			ResponseErrorWithData(c, CodeUserNotExist, code.ErrorSNNotExit)
			return
		}
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
		global.Log.Errorf("controller UpdateUserPwdHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*pkg.MyClaims)
	var data string
	data, err = logic.UpdateUserPwd(&puup, claims.SN) // claims.SN is now int64
	if err != nil {
		global.Log.Errorf("controller UpdateUserPwdHandler logic.UpdateUserPwd err:%s\n", err.Error())
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
		global.Log.Errorf("controller LogoutHandler logic.Logout err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, "注销成功")
}

// DeleteUserListHandler admin 删除用户
func DeleteUserListHandler(c *gin.Context) {
	var pdl model.ParamDeleteList
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteUserListHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.DeleteUserList(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteUserListHandler logic.DeleteUserList err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
