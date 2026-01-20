package message

import (
	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func CreateMessage(msg *model.MessageModel) (string, error) {
	if err := global.DB.Create(msg).Error; err != nil {
		global.Log.Errorf("message CreateMessage err:%s\n", err.Error())
		return "", err
	}
	return code.StrCreateSucceed, nil
}

func CreateMessagesBatch(msgs []model.MessageModel) (string, error) {
	if len(msgs) == 0 {
		return code.StrCreateSucceed, nil
	}
	if err := global.DB.CreateInBatches(msgs, 100).Error; err != nil {
		global.Log.Errorf("message CreateMessagesBatch err:%s\n", err.Error())
		return "", err
	}
	return code.StrCreateSucceed, nil
}

func UpdateMessageRead(sn int64, revUserSN int64) (string, error) {
	result := global.DB.Model(&model.MessageModel{}).
		Where("sn = ? AND rev_user_sn = ?", sn, revUserSN).
		Update("is_read", true)
	if result.Error != nil {
		global.Log.Errorf("message UpdateMessageRead err:%s\n", result.Error.Error())
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", code.ErrorSNNotExist
	}
	return code.StrUpdateSucceed, nil
}

func GetUserMessageList(pl *model.ParamList, userSN int64) ([]model.MessageModel, int64, error) {
	if userSN == 0 {
		return nil, 0, code.ErrorInvalidParam
	}
	return mysql.GetTableList[model.MessageModel]("message_models", pl, "rev_user_sn = ?", userSN)
}

func GetSentMessageList(pl *model.ParamList, userSN int64) ([]model.MessageModel, int64, error) {
	if userSN == 0 {
		return nil, 0, code.ErrorInvalidParam
	}
	return mysql.GetTableList[model.MessageModel]("message_models", pl, "send_user_sn = ?", userSN)
}

func GetAllMessagesList(pl *model.ParamList) ([]model.MessageModel, int64, error) {
	return mysql.GetTableList[model.MessageModel]("message_models", pl, "")
}
