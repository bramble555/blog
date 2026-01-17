package logic

import (
	"time"

	"github.com/bramble555/blog/dao/mysql/banner"
	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"github.com/bramble555/blog/pkg"
)

func PostBindEmail(sn int64, email string) error {
	err := global.DB.Table("user_models").Where("sn = ?", sn).Updates(map[string]interface{}{
		"email": email,
	}).Error

	if err != nil {
		global.Log.Errorf("update email err:%s\n", err.Error())
		return err
	}
	return nil
}

// 用户名或者邮箱登录
func UsernameLogin(peu *model.ParamUsername) (model.ResponseLogin, error) {
	// 判断用户名是否存在
	ok, err := user.CheckUserExistByName(peu.Username)
	resp := model.ResponseLogin{}
	if err != nil {
		return resp, err
	}
	if !ok {
		return resp, code.ErrorUserNotExit
	}
	// 用户存在,判断密码是否错误
	ok, err = user.QueryPasswordByUsername(peu)
	if err != nil {
		return resp, err
	}
	if !ok {
		return resp, code.ErrorPasswordWrong
	}
	return user.GetUserDetail(peu)
}
func RegisterUser(pr *model.ParamRegister) (model.ResponseLogin, error) {
	resp := model.ResponseLogin{}
	role := int64(ctype.PermissionUser)
	err := user.CreateUser(role, pr.Username, pr.Password)
	if err != nil {
		return resp, err
	}
	pu := model.ParamUsername{
		Username: pr.Username,
		Password: pr.Password,
	}
	resp, err = user.GetUserDetail(&pu)
	if err != nil {
		return resp, err
	}
	if pr.Email != "" {
		if err = PostBindEmail(resp.SN, pr.Email); err != nil {
			return resp, err
		}
	}
	return resp, nil
}
func GetUserList(role int64, pl *model.ParamList) (*[]model.UserModel, error) {
	udl, err := user.GetUserList(pl)
	if err != nil {
		return nil, err
	}

	// 手机号和邮箱脱敏处理
	for i := range *udl {
		(*udl)[i].Email = pkg.DesensitizeEmail((*udl)[i].Email)
		(*udl)[i].Phone = pkg.DesensitizePhone((*udl)[i].Phone)
	}

	// 如果是普通用户，username 返回 "****"
	if role == int64(ctype.PermissionUser) {
		for i := range *udl {
			(*udl)[i].Username = "****"
		}
	}

	return udl, nil
}
func UpdateUserRole(puur *model.ParamUpdateUserRole) (string, error) {
	ok, err := user.CheckUserExistBySN(puur.UserSN)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExit
	}
	return user.UpdateUserRole(puur)
}
func UpdateUserPwd(puup *model.ParamUpdateUserPwd, sn int64) (string, error) {
	ok, err := user.CheckPwdExistBySN(sn, puup.OldPwd)
	if err != nil || !ok {
		return "你输入的密码有误,请重新尝试", err
	}
	data, err := user.UpdateUserPwd(puup, sn)
	if err != nil {
		return "", err
	}
	return data, nil
}

func SelectUserBanner(userSN int64, bannerSN int64) (string, error) {
	ok, err := user.CheckUserExistBySN(userSN)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorUserNotExit
	}

	if bannerSN <= 0 {
		return "", code.ErrorSNNotExit
	}

	bd, err := banner.GetBannerBySN(&bannerSN)
	if err != nil {
		return "", err
	}
	avatarPath := "/uploads/file/" + bd.Name
	return user.UpdateUserAvatar(userSN, avatarPath)
}

// Logout 针对注销的操作
func Logout(token string, diff time.Duration) error {
	return redis.Logout(token, diff)
}
func DeleteUser(psn *model.ParamSN) (string, error) {
	ok, err := user.CheckUserExistBySN(psn.SN)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", code.ErrorSNNotExit
	}
	return user.DeleteUser(psn.SN)
}
