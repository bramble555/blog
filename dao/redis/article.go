package redis

import (
	"strconv"
	"time"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg/convert"
	"github.com/go-redis/redis"
)

// SyncLatestArticlesToRedis 从 MySQL 全量同步最新文章到 Redis
// 查询最新的 10 篇文章，清空 Redis 后批量写入
func SyncLatestArticlesToRedis() error {
	// 查询最新的 10 篇文章（按创建时间降序）
	var articles []model.ArticleModel
	err := global.DB.
		Select("sn, create_time").
		Order("create_time DESC").
		Limit(10).
		Find(&articles).Error
	if err != nil {
		global.Log.Errorf("SyncLatestArticlesToRedis: query articles failed: %v", err)
		return err
	}

	key := getKeyName(KeyZSetHomeLatestArticleSN)

	// 使用 Pipeline 提升性能
	pipe := global.Redis.Pipeline()

	// 先清空旧数据
	pipe.Del(key)

	// 批量添加文章，使用创建时间的 Unix 时间戳作为 score
	if len(articles) > 0 {
		members := make([]redis.Z, 0, len(articles))
		for _, article := range articles {
			members = append(members, redis.Z{
				Score:  float64(article.CreateTime.Unix()),
				Member: strconv.FormatInt(article.SN, 10),
			})
		}
		pipe.ZAdd(key, members...)
	}

	// 设置过期时间（10 分钟，防止 Redis 内存泄漏）
	pipe.Expire(key, 10*time.Minute)

	// 执行 Pipeline
	_, err = pipe.Exec()
	if err != nil {
		global.Log.Errorf("SyncLatestArticlesToRedis: pipeline exec failed: %v", err)
		return err
	}

	global.Log.Infof("SyncLatestArticlesToRedis: synced %d articles to Redis", len(articles))
	return nil
}

// AddArticleToLatest 添加单篇文章到 Redis ZSet，并自动淘汰旧数据
// sn: 文章 SN
// timestamp: 文章创建时间戳（Unix 秒）
func AddArticleToLatest(sn int64, timestamp int64) error {
	key := getKeyName(KeyZSetHomeLatestArticleSN)

	pipe := global.Redis.Pipeline()

	// 添加新文章
	pipe.ZAdd(key, redis.Z{
		Score:  float64(timestamp),
		Member: strconv.FormatInt(sn, 10),
	})

	// 只保留最新的 10 篇文章（删除排名 10 之后的数据）
	// ZSet 默认升序，索引 0 是最早的，-11 表示倒数第 11 个（保留 0 到 -11，即最新 10 个）
	pipe.ZRemRangeByRank(key, 0, -11)

	// 刷新过期时间
	pipe.Expire(key, 10*time.Minute)

	_, err := pipe.Exec()
	if err != nil {
		global.Log.Warnf("AddArticleToLatest: failed to add article %d: %v", sn, err)
		return err
	}

	global.Log.Infof("AddArticleToLatest: added article %d to Redis cache", sn)
	return nil
}

// RemoveArticleFromLatest 从 Redis ZSet 中删除指定文章
func RemoveArticleFromLatest(sn int64) error {
	key := getKeyName(KeyZSetHomeLatestArticleSN)

	err := global.Redis.ZRem(key, strconv.FormatInt(sn, 10)).Err()
	if err != nil {
		global.Log.Warnf("RemoveArticleFromLatest: failed to remove article %d: %v", sn, err)
		return err
	}

	global.Log.Infof("RemoveArticleFromLatest: removed article %d from Redis cache", sn)
	return nil
}

// GetLatestArticleSNs 从 Redis 获取最新文章的 SN 列表（倒序，最新的在前）
// limit: 获取数量，默认 10
func GetLatestArticleSNs(limit int) ([]int64, error) {
	key := getKeyName(KeyZSetHomeLatestArticleSN)

	// 使用 ZREVRANGE 获取最新的文章（按 score 倒序）
	// 索引 0 是最新的，limit-1 是第 limit 个
	result, err := global.Redis.ZRevRange(key, 0, int64(limit-1)).Result()
	if err != nil {
		global.Log.Errorf("GetLatestArticleSNs: failed to get SNs from Redis: %v", err)
		return nil, err
	}

	// 转换字符串为 int64
	sns, err := convert.StringSliceToInt64Slice(result)
	if err != nil {
		global.Log.Errorf("GetLatestArticleSNs: failed to convert SNs: %v", err)
		return nil, err
	}

	global.Log.Infof("GetLatestArticleSNs: retrieved %d article SNs from Redis", len(sns))
	return sns, nil
}
