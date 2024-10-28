package ctype

import "encoding/json"

type Role uint

const (
	PermissionAdmin       Role = 1 //管理员
	PermissionUser        Role = 2 //普通用户
	PermissionVisitor     Role = 3 //访客
	PermissionDisableUser Role = 4 //被禁用户
)

// MarshalJSON 方法用于将 Role 类型转换为 JSON 格式。
// 它调用 s.String() 方法获取角色的字符串表示，并将其编码为 JSON
func (s Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s Role) String() string {
	var str string
	switch s {
	case PermissionAdmin:
		str = "管理员"
	case PermissionUser:
		str = "用户"
	case PermissionVisitor:
		str = "游客"
	case PermissionDisableUser:
		str = "被禁用的用户"
	default:
		str = "其他"
	}
	return str
}
