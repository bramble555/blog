package model

import (
	"time"

	"github.com/bramble555/blog/model/ctype"
)

type ChatModel struct {
	*MODEL
	NickName string        `json:"nick_name"`
	Avatar   string        `json:"avatar"`
	Content  string        `json:"content" binding:"required"`
	IP       string        `json:"ip"`
	Addr     string        `json:"addr"`
	MsgType  ctype.MsgType `json:"msg_type" binding:"required"`
}
type ParamChatGroup struct {
	NickName string        `json:"nick_name"`
	Avatar   string        `json:"avatar"`
	Content  string        `json:"content" binding:"required"`
	MsgType  ctype.MsgType `json:"msg_type" binding:"required"`
}
type ResponseChatGroup struct {
	ParamChatGroup
	Date        time.Time `json:"date"`
	OnlineCount int       `json:"online_count"`
}
