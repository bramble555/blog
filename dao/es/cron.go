package es

import (
	"github.com/bramble555/blog/global"
	"github.com/robfig/cron/v3"
)

// StartSyncCronJob 启动定时同步任务
func StartSyncCronJob() *cron.Cron {
	c := cron.New()

	// 每天 00:00 执行全量同步
	// Cron 表达式：0 0 * * * (分 时 日 月 周)
	_, err := c.AddFunc("0 0 * * *", func() {
		global.Log.Info("Cron job triggered: starting sync articles to ES")
		if err := SyncArticlesToES(); err != nil {
			global.Log.Errorf("Cron job: sync articles to ES failed: %v", err)
		} else {
			global.Log.Info("Cron job: sync articles to ES completed successfully")
		}
	})

	if err != nil {
		global.Log.Errorf("StartSyncCronJob add cron job failed: %v", err)
		return nil
	}

	c.Start()
	global.Log.Info("ES sync cron job started successfully (runs daily at 00:00)")

	return c
}
