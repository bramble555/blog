package router

import (
	"expvar"
	"net/http"
	"runtime"
	"time"

	"github.com/bramble555/blog/global"
	"github.com/gin-gonic/gin"
)

func InitMetricsRoutes(r *gin.RouterGroup) gin.IRoutes {
	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	// Publish the database connection pool statistics (for GORM).
	expvar.Publish("database", expvar.Func(func() interface{} {
		sqlDB, err := global.DB.DB() // 注意这里是 global.DB，而不是 global.db
		if err != nil {
			return err.Error()
		}
		return sqlDB.Stats() // 返回 *sql.DB 的连接池统计
	}))

	// Publish the current time.
	expvar.Publish("timeNow", expvar.Func(func() interface{} {
		return time.Now().Format("2006-01-02 15:04:05")
	}))

	// Publish Go memory stats
	expvar.Publish("memoryStats", expvar.Func(func() interface{} {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		return map[string]uint64{
			"Alloc":        memStats.Alloc,         // 当前已分配的字节数
			"HeapIdle":     memStats.HeapIdle,      // 空闲堆内存
			"HeapInuse":    memStats.HeapInuse,     // 使用中的堆内存
			"HeapReleased": memStats.HeapReleased,  // 已释放的堆内存
			"HeapSys":      memStats.HeapSys,       // 堆的总内存
			"NumGC":        uint64(memStats.NumGC), // 完成的 GC 循环数
			"LastGC":       memStats.LastGC,        // 最近一次 GC 的时间
			"Lookups":      memStats.Lookups,       // 全局符号表查找次数
			"GCSys":        memStats.GCSys,         // 垃圾回收使用的内存
		}
	}))

	// 直接在 Gin 中处理 /api/debug/vars 路径的请求
	r.GET("/debug/vars", func(c *gin.Context) {
		// 获取 expvar.Func 并调用它，获得实际值
		goroutines := expvar.Get("goroutines").(expvar.Func)()
		database := expvar.Get("database").(expvar.Func)()
		timeNow := expvar.Get("timeNow").(expvar.Func)()
		memoryStats := expvar.Get("memoryStats").(expvar.Func)()

		// 返回 JSON 格式的响应
		c.JSON(http.StatusOK, map[string]interface{}{
			"goroutines":  goroutines,
			"database":    database,
			"timeNow":     timeNow,
			"memoryStats": memoryStats,
		})
	})

	return r
}
