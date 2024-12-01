package logic

import (
	"time"

	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func GetUserLoginData() ([]model.DailyLoginCount, error) {
	return user.GetUserLoginData()
}
func GetDataSum() (*model.DataSumResponse, error) {
	var res model.DataSumResponse
	var count int64

	// 统计用户数量
	err := global.DB.Table("user_models").Count(&count).Error
	if err != nil {
		global.Log.Errorf("user_models select err: %s\n", err.Error())
		return nil, err
	}
	res.UserCount = count

	// 统计文章数量
	err = global.DB.Table("article_models").Count(&count).Error
	if err != nil {
		global.Log.Errorf("article_models select err: %s\n", err.Error())
		return nil, err
	}
	res.ArticleCount = count

	// 统计消息数量
	err = global.DB.Table("message_models").Count(&count).Error
	if err != nil {
		global.Log.Errorf("message_models select err: %s\n", err.Error())
		return nil, err
	}
	res.MessageCount = count

	// 统计聊天群组数量
	err = global.DB.Table("chat_group_models").Count(&count).Error
	if err != nil {
		global.Log.Errorf("chat_group_models select err: %s\n", err.Error())
		return nil, err
	}
	res.ChatGroupCount = count

	// 统计今天登录用户人数
	today := time.Now().Format("2006-01-02") // 获取当天的日期字符串
	err = global.DB.Table("login_models").
		Where("DATE(create_time) = ?", today). 
		Count(&count).Error
	if err != nil {
		global.Log.Errorf("login_models select err: %s\n", err.Error())
		return nil, err
	}
	res.NowLoginCount = count

	return &res, nil
}
