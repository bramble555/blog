package user

import (
	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func GetUserList(pl *model.ParamList) (*[]model.UserModel, error) {
	udl, err := mysql.GetTableList[model.UserModel]("user_models", pl, "")
	if err != nil {
		global.Log.Errorf("user GetUserList err:%s\n", err.Error())
		return nil, err
	}
	return &udl, nil
}
func GetUserDetail(id uint) (*model.UserDetail, error) {
	ud := model.UserDetail{}
	err := global.DB.Table("user_models").Where("id = ?", id).Scan(&ud).Error
	if err != nil {
		global.Log.Errorf("user GetUserDetail err:%s\n", err.Error())
		return nil, err
	}
	return &ud, nil
}
