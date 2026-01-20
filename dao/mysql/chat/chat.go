package chat

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func CreateChatRecord(cm *model.ChatModel) error {
	return global.DB.Create(cm).Error
}

// ListRecentChatRecords 获取最近的聊天记录
func ListRecentChatRecords(limit int) ([]model.ChatModel, error) {
	if limit <= 0 {
		limit = 20
	}
	list := make([]model.ChatModel, 0, limit)
	err := global.DB.Table("chat_models").Order("create_time DESC, sn DESC").Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	//  翻转列表，使最新的记录在最前面,这里手动翻转式为了提高性能
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return list, nil
}

func GetChatRecords(page, size int) ([]model.ChatModel, int64, error) {
	list := make([]model.ChatModel, 0)
	var count int64
	offset := (page - 1) * size

	db := global.DB.Table("chat_models")
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := db.Order("create_time DESC").Limit(size).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
