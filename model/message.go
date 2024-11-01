package model

import "time"

type MessageModel struct {
	ID         uint      `json:"id"`           // 消息 ID
	CreateTime time.Time `json:"create_time"`  // 创建时间
	UpdateTime time.Time `json:"update_time"`  // 更新时间
	SendUserID uint      `json:"send_user_id"` // 发送人 ID
	RevUserID  uint      `json:"rev_user_id"`  // 接收人 ID
	IsRead     bool      `json:"is_read"`      // 接收方是否查看
	Content    string    `json:"content"`      // 消息内容
}
type ParamMessage struct {
	SendUserID uint   `json:"send_user_id"`               // 发送人ID
	RevUserID  uint   `json:"rev_user_id"`                // 接收人ID
	Content    string `json:"content" binding:"required"` // 消息内容
}
type RespondMessage struct {
	MODEL
	SendUserID        uint   `json:"send_user_id"` // 发送人id
	SendUserNicekName string `json:"send_user_nick_name"`
	SendUserAvatar    string `json:"send_user_avater"`
	RevUserID         uint   `json:"rev_user_id"` // 接收人id
	RevUserNicekName  string `json:"rev_user_nick_name"`
	RevUserAvatar     string `json:"rev_user_avater"`
	IsRead            bool   `json:"is_read"` // 接收方是否查看
	Content           string `json:"content"` // 消息内容
}
