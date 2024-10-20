package redis

import (
	"context"
	"time"

	"github.com/bramble555/blog/global"
	"github.com/go-redis/redis"
)

func Init() (*redis.Client, error) {
	redisConf := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr(),
		Password: redisConf.Password,
		DB:       0,
		PoolSize: redisConf.PoolSize,
	})
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
