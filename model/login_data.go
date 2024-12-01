package model

import "github.com/bramble555/blog/model/ctype"

type LoginDataModel struct {
	*MODEL
	UserID    uint             `json:"user_id"`
	IP        string           `json:"ip"` // 登录的 IP
	NickName  string           `json:"nick_name"`
	Token     string           `json:"token"`
	Device    string           `json:"device"` // 登录设备
	Addr      string           `json:"addr"`
	LoginType ctype.SignStatus `json:"login_type"` // 登录方式
}
