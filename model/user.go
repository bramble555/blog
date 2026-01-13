package model

import (
	"time"

	"github.com/bramble555/blog/model/ctype"
)

type UserModel struct {
	MODEL
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
	ArticleSN  int64            `json:"article_sn"` //发布的文章
	CollectSN  int64            `json:"collect_sn"`
}
type ParamFlagUser struct {
	SN       int64  `json:"sn"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int64  `json:"role"`
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
	UserSN int64      `json:"user_sn,string"`
	Role   ctype.Role `json:"role"`
}
type ParamUpdateUserPwd struct {
	OldPwd string `json:"old_pwd"`
	Pwd    string `json:"pwd"`
}
type UserDetail struct {
	SN       int64  `json:"user_sn,string"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"-"`
}
type DailyLoginCount struct {
	LoginDate  time.Time `json:"login_date"`  // 登录日期
	LoginCount int64     `json:"login_count"` // 登录次数
}

type ResponseLogin struct {
	Token    string `json:"token"`
	SN       int64  `json:"sn,string"` // Snowflake ID as string
	Username string `json:"username"`
	Role     int64  `json:"role"`
}
