package model

type MessageModel struct {
	*MODEL
	SendUserID     uint   `json:"send_user_id,string"` // 发送人 ID
	SendUsername   string `json:"send_user_name"`
	SendUserAvatar string `json:"send_user_avater"`
	RevUserID      uint   `json:"rev_user_id,string"` // 接收人 ID
	RevUsername    string `json:"rev_user_name"`
	RevUserAvatar  string `json:"rev_user_avater"`
	IsRead         bool   `json:"is_read"`                    // 接收方是否查看
	Content        string `json:"content" binding:"required"` // 消息内容
}
type ParamMessage struct {
	SendUserID uint   `json:"send_user_id,string"`        // 发送人ID
	RevUserID  uint   `json:"rev_user_id,string"`         // 接收人ID
	Content    string `json:"content" binding:"required"` // 消息内容
}
type RespondMessage struct {
	*MODEL
	SendUserID     uint   `json:"send_user_id,string"` // 发送人id
	SendUsername   string `json:"send_user_name"`
	SendUserAvatar string `json:"send_user_avater"`
	RevUserID      uint   `json:"rev_user_id,string"` // 接收人id
	RevUsername    string `json:"rev_user_name"`
	RevUserAvatar  string `json:"rev_user_avater"`
	IsRead         bool   `json:"is_read"`       // 接收方是否查看
	Content        string `json:"content"`       // 消息内容
	MessageCount   uint   `json:"message_count"` // 消息数量
}
type ParamRecordID struct {
	UserID uint `json:"user_id,string"` // 查询聊天记录的 ID
}
type RespondMessageGroup map[uint]*RespondMessage
