package code

import "errors"

var (
	ErrorIDExit        = errors.New("ID 已存在")
	
	ErrorIDNotExit     = errors.New("ID 不存在")
	ErrorUserExit      = errors.New("用户已存在")
	ErrorUserNotExit   = errors.New("用户不存在")
	ErrorPasswordWrong = errors.New("密码错误")
	ErrorTitleExit     = errors.New("主题已存在") // 广告主题
)
var (
	CreateSucceed = "创建成功"
)
