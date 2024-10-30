package pkg

import "strings"

// 邮箱脱敏
func DesensitizeEmail(email string) string {
	// x*******@qq.com
	eList := strings.Split(email, "@")
	if len(eList) != 2 {
		return ""
	}
	return eList[0][:1] + "****" + eList[1]
}

// 手机号脱敏
func DesensitizePhone(tel string) string {
	// 139****0000
	// 不是国内手机号直接返回 ""
	if len(tel) != 11 {
		return ""
	}

	return tel[:3] + "****" + tel[7:]
}
