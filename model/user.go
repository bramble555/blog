package model

import "github.com/bramble555/blog/model/ctype"

type UserModel struct {
	*MODEL
	Username   string           `json:"username"`
	PassWord   string           `json:"-"`
	Avatar     string           `json:"avatar"`
	Email      string           `json:"email"`
	Tel        string           `json:"tel"`
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
type ParamEmailUser struct {
	Username string `json:"username" binding:"required" msg:"请输入用户名"`
	Password string `json:"password" binding:"required" msg:"请输入密码"`
}

// 返回响应
type ResponseUserLogin struct {
	Token    string
	UserName string
	Role     int
}
