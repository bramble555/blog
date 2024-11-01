package redis

import (
	"errors"
	"time"

	"github.com/bramble555/blog/global"
	"github.com/go-redis/redis"
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
	_, err := global.Redis.Get(getKeyName(token)).Result()
	if err == redis.Nil {
		// 键不存在
		return false
	} else if err != nil {
		// 其他错误
		global.Log.Errorf("redis CheckLogout err: %s\n", err.Error())
		return false
	}
	// 存在
	return true
}
