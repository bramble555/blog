package redis

import (
	"fmt"
	"strconv"

	"github.com/bramble555/blog/global"
	"gorm.io/gorm"
)

// Redis Hash Fields for Article & Comment Counts
const (
	FieldLook        = "look"
	FieldDigg        = "digg"
	FieldCommentDigg = "comment_digg" // 评论的点赞数
)

// IncrArticleCount 增加或减少文章计数器 (大厂高并发推荐方案)
// sn: 文章 SN
// field: 计数维度 (look/digg/comment_digg)
// delta: 增量 (1 或 -1)
func IncrArticleCount(sn int64, field string, delta int64) error {
	var key string
	switch field {
	case FieldLook:
		key = getKeyName(KeyHashArticleLookCount)
	case FieldDigg:
		key = getKeyName(KeyHashArticleDiggCount)
	case FieldCommentDigg:
		key = getKeyName(KeyHashCommentDiggCount)
	default:
		return nil
	}

	// HINCRBY 保证原子性
	err := global.Redis.HIncrBy(key, strconv.FormatInt(sn, 10), delta).Err()
	if err != nil {
		global.Log.Errorf("IncrArticleCount: HIncrBy failed: key=%s, sn=%d, err=%v", key, sn, err)
		return err
	}
	return nil
}

// GetRedisArticlesCounts 批量获取文章计数
// sns: 文章 SN 列表
// field: 计数维度 (look/digg/comment_digg)
// 返回: map[sn]count
func GetRedisArticlesCounts(sns []int64, field string) (map[int64]int64, error) {
	if len(sns) == 0 {
		return make(map[int64]int64), nil
	}

	var key string
	switch field {
	case FieldLook:
		key = getKeyName(KeyHashArticleLookCount)
	case FieldDigg:
		key = getKeyName(KeyHashArticleDiggCount)
	case FieldCommentDigg:
		key = getKeyName(KeyHashCommentDiggCount)
	default:
		return make(map[int64]int64), nil
	}

	// 将 int64 SN 转换为 string field
	fields := make([]string, len(sns))
	for i, sn := range sns {
		fields[i] = strconv.FormatInt(sn, 10)
	}

	// HMGet 批量获取
	results, err := global.Redis.HMGet(key, fields...).Result()
	if err != nil {
		global.Log.Errorf("GetArticlesCounts: HMGet failed: key=%s, err=%v", key, err)
		return nil, err
	}

	countMap := make(map[int64]int64, len(sns))
	for i, result := range results {
		if result == nil {
			countMap[sns[i]] = 0
			continue
		}
		// Redis 返回的是 string 类型
		valStr, ok := result.(string)
		if !ok {
			countMap[sns[i]] = 0
			continue
		}
		val, _ := strconv.ParseInt(valStr, 10, 64)
		countMap[sns[i]] = val
	}

	return countMap, nil
}

// GetRedisAllCounts 获取某个维度的所有缓存计数 (用于定时同步 MySQL)
func GetRedisAllCounts(field string) (map[int64]int64, error) {
	var key string
	switch field {
	case FieldLook:
		key = getKeyName(KeyHashArticleLookCount)
	case FieldDigg:
		key = getKeyName(KeyHashArticleDiggCount)
	case FieldCommentDigg:
		key = getKeyName(KeyHashCommentDiggCount)
	default:
		return nil, nil
	}

	result, err := global.Redis.HGetAll(key).Result()
	if err != nil {
		global.Log.Errorf("GetAllCounts: HGetAll failed: key=%s, err=%v", key, err)
		return nil, err
	}

	counts := make(map[int64]int64, len(result))
	for snStr, countStr := range result {
		sn, err := strconv.ParseInt(snStr, 10, 64)
		if err != nil {
			continue
		}
		count, _ := strconv.ParseInt(countStr, 10, 64)
		counts[sn] = count
	}

	return counts, nil
}

// SyncCountsToMySQL 将 Redis 中的实时计数批量同步回 MySQL (大厂标准：定时异步落库)
func SyncCountsToMySQL() error {
	fields := []string{FieldLook, FieldDigg, FieldCommentDigg}

	for _, field := range fields {
		counts, err := GetRedisAllCounts(field)
		if err != nil {
			continue
		}
		if len(counts) == 0 {
			continue
		}

		global.Log.Infof("SyncCountsToMySQL: syncing %s counts for %d items", field, len(counts))

		for sn, count := range counts {
			if count == 0 {
				continue
			}
			var col string
			var table string = "article_models"

			switch field {
			case FieldLook:
				col = "look_count"
			case FieldDigg:
				col = "digg_count"
			case FieldCommentDigg:
				col = "digg_count"
				table = "comment_models"
			}

			// 1. 写回 MySQL：使用 SQL 表达式实现原子叠加 (look_count = look_count + delta)
			err := global.DB.Table(table).
				Where("sn = ?", sn).
				UpdateColumn(col, gorm.Expr(fmt.Sprintf("%s + ?", col), count)).Error
			if err != nil {
				global.Log.Errorf("SyncCountsToMySQL: failed to update SN %d, col %s in %s: %v", sn, col, table, err)
				continue
			}

			// 2. 重置 Redis：减去已同步的增量 (不能直接 DEL，防止丢失同步期间产生的新增量)
			_ = IncrArticleCount(sn, field, -count)
		}
	}

	global.Log.Info("SyncCountsToMySQL: sync finished.")
	return nil
}
