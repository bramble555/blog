package redis

import (
	"errors"
	"time"

	"github.com/bramble555/blog/global"
)

func Logout(token string, diff time.Duration) error {
	if diff <= 0 {
		global.Log.Warnf("Invalid expiration time for token logout: %v", diff)
		return errors.New("invalid expiration time")
	}

	err := global.Redis.Set(getKeyName(token), "", diff).Err()
	if err != nil {
		global.Log.Errorf("redis Logout err: %s\n", err.Error())
		return err
	}
	return nil
}

// CheckLogout 检查 token 是否存在于用户注销的 redis 里面
func CheckLogout(token string) bool {
	keys := global.Redis.Keys(getKeyName(token)).Val()
	for _, v := range keys {
		if v == token {
			return true
		}
	}
	return false
}
