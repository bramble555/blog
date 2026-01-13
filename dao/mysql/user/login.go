package user

import (
	"errors"
	"sort"
	"time"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
)

// CheckUserExistByName 检查用户是否存在
func CheckUserExistByName(name string) (bool, error) {
	var count int64
	err := global.DB.Table("user_models").Where("username = ?", name).Count(&count).Error
	if err != nil {
		global.Log.Errorf("user CheckUserExistByName err:%s\n", err.Error())
		return false, err
	}
	return count == 1, nil
}
func CheckUserExistBySN(sn int64) (bool, error) {
	var count int64
	err := global.DB.Table("user_models").Where("sn = ?", sn).Count(&count).Error
	if err != nil {
		global.Log.Errorf("user CheckUserExistBySN err:%s\n", err.Error())
		return false, err
	}
	return count == 1, nil
}

// CheckPwdExistBySN 传入 SN，检查密码是否正确
func CheckPwdExistBySN(sn int64, pwd string) (bool, error) {
	var encryPassword string
	err := global.DB.Table("user_models").Where("sn = ?", sn).
		Select("password").Scan(&encryPassword).Error
	if err != nil {
		global.Log.Errorf("user QueryPassword err: %v\n", err)
		return false, err
	}

	// 如果密码为空，用户不存在
	if encryPassword == "" {
		return false, errors.New("用户不存在")
	}

	// 比较密码
	err = pkg.ComparePasswords(encryPassword, pwd)
	if err != nil {
		return false, errors.New("密码不正确")
	}
	return true, nil
}

// QueryPasswordByUsername 传入username 和 密码 检验密码是否正确，实现登录功能
func QueryPasswordByUsername(peu *model.ParamUsername) (bool, error) {
	var encryPassword string
	err := global.DB.Table("user_models").Where("username = ?", peu.Username).
		Select("password").Scan(&encryPassword).Error // 使用 Scan 将结果绑定到 password
	if err != nil {
		global.Log.Errorf("user QueryPassword err: %v\n", err)
		return false, err
	}
	err = pkg.ComparePasswords(encryPassword, peu.Password)
	if err != nil {
		global.Log.Errorf("user pkg.ComparePassword serr: %v\n", err)
		return false, code.ErrorPasswordWrong
	}
	return true, nil
}
func GetToken(peu *model.ParamUsername) (model.ResponseLogin, error) {
	type paramUserDetail struct {
		SN       int64 // Snowflake ID
		Username string
		Role     int64
	}
	var udd paramUserDetail

	err := global.DB.Table("user_models").Where("username = ?", peu.Username).
		Select("sn,username,role").Scan(&udd).Error
	if err != nil {
		global.Log.Errorf("user GetToken select err:%s\n", err.Error())
		return model.ResponseLogin{}, err
	}
	token, err := pkg.GenToken(udd.SN, udd.Role, udd.Username)
	if err != nil {
		global.Log.Errorf("pkg GetToken err:%s\n", err.Error())
		return model.ResponseLogin{}, err
	}
	return model.ResponseLogin{
		Token:    token,
		SN:       udd.SN, // Return Snowflake ID
		Username: udd.Username,
		Role:     udd.Role,
	}, nil
}
func PostLogin(username string) error {
	type data struct {
		Username string
	}
	d := data{}
	d.Username = username
	err := global.DB.Table("login_models").Create(&d).Error
	if err != nil {
		global.Log.Errorf("create login_models err:%s\n", err.Error())
	}
	return nil
}

func GetUserLoginData() ([]model.DailyLoginCount, error) {
	var queryResults []model.DailyLoginCount

	// 查询数据库，获取过去 7 天有记录的日期和登录次数
	err := global.DB.Table("login_models").
		Select("DATE(created_at) AS login_date, COUNT(*) AS login_count").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Group("DATE(created_at)").
		Order("login_date").
		Scan(&queryResults).Error

	if err != nil {
		global.Log.Errorf("查询 login_models 失败 err:%s\n", err.Error())
		return nil, err
	}

	// 生成过去 7 天的完整日期范围
	dates := make(map[string]int)
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dates[date] = 0
	}

	// 将查询结果填充到日期范围中
	for _, result := range queryResults {
		dateStr := result.LoginDate.Format("2006-01-02") // 转换时间为字符串
		if _, exists := dates[dateStr]; exists {
			dates[dateStr] = int(result.LoginCount)
		}
	}

	// 转换回结果数组
	var finalResults []model.DailyLoginCount
	for dateStr, count := range dates {
		parsedDate, _ := time.Parse("2006-01-02", dateStr) // 转换字符串为 time.Time
		finalResults = append(finalResults, model.DailyLoginCount{
			LoginDate:  parsedDate,
			LoginCount: int64(count),
		})
	}

	// 按日期排序
	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].LoginDate.Before(finalResults[j].LoginDate)
	})

	return finalResults, nil
}
