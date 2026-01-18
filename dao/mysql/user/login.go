package user

import (
	"errors"
	"sort"
	"time"

	"github.com/bramble555/blog/dao/mysql/code"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"gorm.io/gorm"
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

// QueryPasswordByUsername 传入 username 和 密码 检验密码是否正确，实现登录功能
func QueryPasswordByUsername(peu *model.ParamUsername) (bool, error) {
	var encryPassword string
	err := global.DB.Table("user_models").Where("username = ?", peu.Username).
		Select("password").Scan(&encryPassword).Error // 使用 Scan 将结果绑定到 encryPassword
	if err != nil {
		global.Log.Errorf("user QueryPassword err: %v\n", err)
		return false, err
	}
	// 对比密码是否一致
	err = pkg.ComparePasswords(encryPassword, peu.Password)
	if err != nil {
		global.Log.Errorf("user pkg.ComparePassword serr: %v\n", err)
		return false, code.ErrorPasswordWrong
	}
	return true, nil
}

func GetUserDetail(peu *model.ParamUsername) (model.ResponseLogin, error) {
	// 内存对齐
	type paramUserDetail struct {
		Username string
		SN       int64
		Role     int64
		Avatar   string
	}

	// 使用 var 块统一声明所有变量
	var (
		pud   paramUserDetail
		res   model.ResponseLogin // 预声明返回对象，初始即为零值 {}
		token string
		err   error
	)

	// 1. 查询数据库
	// 使用 First/Token 进行查询是否存在,如果不存在 error 是 gorm.ErrRecordNotFound, 不用 Find/Scan
	err = global.DB.Table("user_models").
		Select("sn", "username", "role", "avatar").
		Where("username = ?", peu.Username).
		Take(&pud).Error

	// 细分 err
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warnf("GetToken 失败: 用户 [%s] 不存在", peu.Username)
			return res, errors.New("用户不存在") // 直接返回 res，无需写 model.ResponseLogin{}
		}
		global.Log.Errorf("GetToken 数据库查询异常: %v", err)
		return res, err
	}

	// 2. 生成 Token
	token, err = pkg.GenToken(pud.SN, pud.Role, pud.Username)
	if err != nil {
		global.Log.Errorf("GenToken 签名失败: %v", err)
		return res, err
	}

	// 3. 填充预声明的对象
	res.Token = token
	res.SN = pud.SN
	res.Username = pud.Username
	res.Role = pud.Role
	res.Avatar = pud.Avatar

	return res, nil
}

func PostLogin(username string) error {
	var lm model.LoginModel
	lm.Username = username
	err := global.DB.Table("login_models").Create(&lm).Error
	if err != nil {
		global.Log.Errorf("create login_models err:%s\n", err.Error())
		return err
	}
	return nil
}

func GetUserLoginData() ([]model.DailyLoginCount, error) {
	var queryResults []model.DailyLoginCount

	// 查询数据库，获取过去 7 天有记录的日期和登录次数
	err := global.DB.Table("login_models").
		Select("DATE(create_time) AS login_date, COUNT(*) AS login_count").
		Where("create_time >= ?", time.Now().AddDate(0, 0, -7)).
		Group("DATE(create_time)").
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
