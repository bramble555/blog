package code

import "errors"

var (
	ErrorInvalidParam        = errors.New("请求参数错误")
	ErrorUserExist           = errors.New("用户已存在")
	ErrorUserNotExist        = errors.New("用户不存在")
	ErrorInvalidPassword     = errors.New("用户名或者密码错误")
	ErrorInvalidVerification = errors.New("验证码错误")
	ErrorServerBusy          = errors.New("服务器繁忙")
	ErrorNeedLogin           = errors.New("用户请登录")
	ErrorInvalidAuth         = errors.New("token无效")
	ErrorInvalidID           = errors.New("无效ID")
	ErrorTitleExist          = errors.New("主题已存在")
	ErrorTagNotExist         = errors.New("请输入已有的tag")
	ErrorTagExist            = errors.New("tag 已存在")
)

var (
	ErrorSNNotExist    = ErrorInvalidID
	ErrorPasswordWrong = ErrorInvalidPassword
)

var (
	StrCreateSucceed = "创建成功"
	StrUpdateSucceed = "更新成功"
)
