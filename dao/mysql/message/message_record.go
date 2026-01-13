package message

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func MessageRecord(userSN int64, revUserSN int64) ([]model.MessageModel, error) {
	var messageList []model.MessageModel
	// 查找所有与自己和对方相关的消息
	err := global.DB.Order("create_time asc").
		Find(&messageList, "send_user_sn = ? and rev_user_sn = ? or send_user_sn = ? and rev_user_sn = ?",
			userSN, revUserSN, revUserSN, userSN).Error
	if err != nil {
		return nil, err
	}
	// 标记对方发给自己的消息为已读
	var sns []int64
	for _, m := range messageList {
		// 只有接收者是自己的时候，才标记为已读
		if m.RevUserSN == userSN {
			sns = append(sns, m.SN)
		}
	}
	// 批量更新已读状态
	if len(sns) > 0 {
		global.DB.Model(&model.MessageModel{}).Where("sn IN ?", sns).Update("is_read", true)
	}

	return messageList, nil
}
