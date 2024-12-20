package model

import (
	"time"

	"github.com/bramble555/blog/model/ctype"
)

type UserModel struct {
	*MODEL
	Username   string           `json:"username"`
	PassWord   string           `json:"-"`
	Avatar     string           `json:"avatar"`
	Email      string           `json:"email"`
	Phone      string           `json:"phone"`
	Addr       string           `json:"addr"`
	Token      string           `json:"token"`
	IP         string           `json:"ip"`
	Role       ctype.Role       `json:"role"` //角色权限
	SignStatus ctype.SignStatus `json:"sign_status"`
	ArticleID  uint             `json:"artcile_id"` //发布的文章
	CollectID  uint             `json:"collect_id"`
}
type ParamFlagUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     uint   `json:"role"`
	Avatar   string `json:"avatar"`
}
type ParamUsername struct {
	Username string `json:"username" binding:"required" msg:"请输入用户名"`
	Password string `json:"password" binding:"required" msg:"请输入密码"`
}
type ParamBindEmail struct {
	Email string  `json:"email" binding:"required" msg:"邮箱非法"`
	Code  *string `json:"code"`
}
type ParamUpdateUserRole struct {
	UserID uint       `json:"user_id,string"`
	Role   ctype.Role `json:"role"`
}
type ParamUpdateUserPwd struct {
	OldPwd string `json:"old_pwd"`
	Pwd    string `json:"pwd"`
}
type UserDetail struct {
	ID       uint   `json:"user_id,string"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"-"`
}
type DailyLoginCount struct {
	LoginDate  time.Time `json:"login_date"`  // 登录日期
	LoginCount int       `json:"login_count"` // 登录次数
}

// 返回响应
// type ResponseUserLogin struct {
// 	Token    string
// 	UserName string
// 	Role     int
// }
