package model

type MessageModel struct {
	MODEL
	SendUserSN     int64  `json:"send_user_sn,string"` // 发送人 SN
	SendUsername   string `json:"send_user_name"`
	SendUserAvatar string `json:"send_user_avater"`
	RevUserSN      int64  `json:"rev_user_sn,string"` // 接收人 SN
	RevUsername    string `json:"rev_user_name"`
	RevUserAvatar  string `json:"rev_user_avater"`
	IsRead         bool   `json:"is_read"`                    // 接收方是否查看
	Content        string `json:"content" binding:"required"` // 消息内容
}
type ParamMessage struct {
	SendUserSN int64  `json:"send_user_sn,string"`        // 发送人SN
	RevUserSN  int64  `json:"rev_user_sn,string"`         // 接收人SN
	Content    string `json:"content" binding:"required"` // 消息内容
}
type RespondMessage struct {
	MODEL
	SendUserSN     int64  `json:"send_user_sn,string"` // 发送人sn
	SendUsername   string `json:"send_user_name"`
	SendUserAvatar string `json:"send_user_avater"`
	RevUserSN      int64  `json:"rev_user_sn,string"` // 接收人sn
	RevUsername    string `json:"rev_user_name"`
	RevUserAvatar  string `json:"rev_user_avater"`
	IsRead         bool   `json:"is_read"`       // 接收方是否查看
	Content        string `json:"content"`       // 消息内容
	MessageCount   int64  `json:"message_count"` // 消息数量
}
type ParamRecordSN struct {
	UserSN int64 `json:"user_sn,string"` // 查询聊天记录的 SN
}
type RespondMessageGroup map[int64]*RespondMessage
