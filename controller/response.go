package controller

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/gin-gonic/gin"
)

type resCode int

const (
	CodeSucceed = 10000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeInvalidVerification
	CodeServerBusy // 连接不上数据库

	CodeNeedLogin
	CodeInvalidAuth

	CodeInvalidID
	CodeTitleExist
	CodeTagNotExist
	CodeTagExist
	CodeTooManyRequests
)

var codeMsgMap = map[resCode]string{
	CodeSucceed:             "succeed",
	CodeInvalidParam:        "请求参数错误",
	CodeUserExist:           "用户已存在",
	CodeUserNotExist:        "用户不存在",
	CodeInvalidPassword:     "用户名或者密码错误",
	CodeInvalidVerification: "验证码错误",
	CodeServerBusy:          "服务器繁忙",
	CodeNeedLogin:           "用户请登录",
	CodeInvalidAuth:         "token无效",
	CodeInvalidID:           "无效ID",
	CodeTitleExist:          "主题已存在",
	CodeTagNotExist:         "请输入已有的tag",
	CodeTagExist:            "tag 已存在",
	CodeTooManyRequests:     "请求过于频繁",
}

type responseData struct {
	Code resCode `json:"code"`
	Msg  any     `json:"msg"`
	Data any     `json:"data"`
}

func (rc resCode) msg() string {
	msg, ok := codeMsgMap[rc]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}

func codeFromError(err error) resCode {
	if err == nil {
		return CodeSucceed
	}
	switch {
	case errors.Is(err, code.ErrorInvalidParam):
		return CodeInvalidParam
	case errors.Is(err, code.ErrorUserExist):
		return CodeUserExist
	case errors.Is(err, code.ErrorUserNotExist):
		return CodeUserNotExist
	case errors.Is(err, code.ErrorInvalidPassword):
		return CodeInvalidPassword
	case errors.Is(err, code.ErrorInvalidVerification):
		return CodeInvalidVerification
	case errors.Is(err, code.ErrorServerBusy):
		return CodeServerBusy
	case errors.Is(err, code.ErrorNeedLogin):
		return CodeNeedLogin
	case errors.Is(err, code.ErrorInvalidAuth):
		return CodeInvalidAuth
	case errors.Is(err, code.ErrorInvalidID):
		return CodeInvalidID
	case errors.Is(err, code.ErrorTitleExist):
		return CodeTitleExist
	case errors.Is(err, code.ErrorTagNotExist):
		return CodeTagNotExist
	case errors.Is(err, code.ErrorTagExist):
		return CodeTagExist
	default:
		return CodeServerBusy
	}
}

// ResponseSucceed 成功响应
func ResponseSucceed(c *gin.Context, data any) {
	c.JSON(200, &responseData{
		Code: CodeSucceed,
		Msg:  codeMsgMap[CodeSucceed],
		Data: data,
	})
}

// ResponseError 返回错误，但是不知道啥错误，所以要传入code
func ResponseError(c *gin.Context, code resCode) {
	c.JSON(200, &responseData{
		Code: code,
		Msg:  code.msg(),
		Data: "",
	})
}

func ResponseErrorByErr(c *gin.Context, err error) {
	ResponseError(c, codeFromError(err))
}

// 返回错误附带数据
func ResponseErrorWithData(c *gin.Context, code resCode, data any) {
	c.JSON(200, &responseData{
		Code: code,
		Msg:  code.msg(),
		Data: data,
	})
}

func ResponseErrorWithErr(c *gin.Context, err error, data any) {
	ResponseErrorWithData(c, codeFromError(err), data)
}
