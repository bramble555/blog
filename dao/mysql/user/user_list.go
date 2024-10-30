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
