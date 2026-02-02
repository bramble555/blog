package redis

import (
	"github.com/bramble555/blog/global"
	"github.com/robfig/cron/v3"
)

// StartSyncCronJob 启动 Redis 缓存定时同步任务
// 每 5 分钟从 MySQL 同步最新文章到 Redis
func StartSyncCronJob() *cron.Cron {
	c := cron.New()

	// 每 5 分钟执行一次全量同步
	// Cron 表达式：*/5 * * * * (每 5 分钟)
	_, err := c.AddFunc("*/5 * * * *", func() {
		global.Log.Info("Cron job triggered: starting sync latest articles to Redis")
		if err := SyncLatestArticlesToRedis(); err != nil {
			global.Log.Errorf("Cron job: sync latest articles to Redis failed: %v", err)
		} else {
			global.Log.Info("Cron job: sync latest articles to Redis completed successfully")
		}
	})

	if err != nil {
		global.Log.Errorf("StartSyncCronJob: add cron job failed: %v", err)
		return nil
	}

	// 2. 每 5 分钟将最新的文章计数器同步回 MySQL
	_, err = c.AddFunc("*/5 * * * *", func() {
		global.Log.Info("Cron Job: starting sync counts to MySQL...")
		if err := SyncCountsToMySQL(); err != nil {
			global.Log.Errorf("Cron Job Error: SyncCountsToMySQL failed: %v", err)
		}
	})
	if err != nil {
		global.Log.Fatalf("Cron Job: failed to add SyncCountsToMySQL task: %v", err)
	}

	c.Start()
	global.Log.Info("Redis Sync Cron Job started.")

	// 立即执行一次全量同步（服务启动时预热缓存）
	go func() {
		global.Log.Info("Initial sync: starting sync latest articles to Redis")
		if err := SyncLatestArticlesToRedis(); err != nil {
			global.Log.Errorf("Initial sync: sync latest articles to Redis failed: %v", err)
		} else {
			global.Log.Info("Initial sync: sync latest articles to Redis completed successfully")
		}
	}()

	return c
}
