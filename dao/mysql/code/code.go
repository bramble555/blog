package code

import "errors"

var (
	ErrorSNExist = errors.New("SN 已存在")

	ErrorSNNotExist    = errors.New("SN 不存在")
	ErrorUserExist     = errors.New("用户已存在")
	ErrorUserNotExist  = errors.New("用户不存在")
	ErrorPasswordWrong = errors.New("密码错误")
	ErrorTitleExist    = errors.New("主题已存在")

	ErrorCreateWrong     = errors.New("创建错误")
	ErrorAssertionFailed = errors.New("断言失败")
	ErrorTagNotExist     = errors.New("tag 不存在")
	ErrorTagExist        = errors.New("tag 已存在")
)
var (
	StrCreateSucceed = "创建成功"
	StrUpdateSucceed = "更新成功"
)
