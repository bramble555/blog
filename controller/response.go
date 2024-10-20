package controller

import "github.com/gin-gonic/gin"

/*
 包装一个函数 可以实现
 {
	"resCode": 10000,
	"msg": ...,
	"data": ...,
 }
*/
type resCode int

const (
	CodeSucceed = 10000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy // 连接不上数据库

	CodeNeedLogin
	CodeInvalidAuth

	CodeInvalidID
)

var codeMsgMap = map[resCode]string{
	CodeSucceed:         "succeed",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或者密码错误",
	CodeServerBusy:      "服务器繁忙",
	CodeNeedLogin:       "用户请登录",
	CodeInvalidAuth:     "token无效",
	CodeInvalidID:       "无效ID",
}

type ResponseData struct {
	Code resCode `json:"code"`
	Msg  any     `json:"msg"`
	Data any     `json:"data"`
}

func (rc resCode) Msg() string {
	msg, ok := codeMsgMap[rc]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}

// ResponseSucceed 成功响应
func ResponseSucceed(c *gin.Context, data any) {
	c.JSON(200, &ResponseData{
		Code: CodeSucceed,
		Msg:  codeMsgMap[CodeSucceed],
		Data: data,
	})
}

// ResponseError 返回错误，但是不知道啥错误，所以要传入code
func ResponseError(c *gin.Context, code resCode) {
	c.JSON(200, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: "",
	})
}

// 返回错误附带数据
func ResponseErrorWithData(c *gin.Context, code resCode, data any) {
	c.JSON(200, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: data,
	})
}
