package main

import (
	"context"
	_ "database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bramble555/blog/dao/es"
	dao_mysql "github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/flag"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logger"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/router"
	"github.com/bramble555/blog/setting"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var err error
	// 初始化 config
	if err = setting.Init(); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
		return
	}
	// 初始化 logger
	if global.Log, err = logger.Init(); err != nil {
		fmt.Printf("Init logger failed, err:%v\n", err)
		return
	}

	// 初始化 db
	if global.DB, err = dao_mysql.Init(); err != nil {
		global.Log.Errorf("Init mysql failed, err:%v\n", err)
		return
	}

	// 迁移数据库
	db, _ := global.DB.DB()
	migrationDriver, err := mysql.WithInstance(db, &mysql.Config{}) // 使用 migrate 的 mysql 驱动
	if err != nil {
		global.Log.Errorf("migration Driver failed, err:%v\n", err)
		return
	}
	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migration",     // 迁移脚本路径
		global.Config.Mysql.DB, // 数据库名
		migrationDriver,        // 迁移驱动
	)
	if err != nil {
		global.Log.Errorf("migration failed, err:%v\n", err)
		return
	}
	err = migrator.Up()
	if err != nil {
		global.Log.Errorf("migration up failed, err:%v\n", err)
		return
	}

	// 初始化 redis
	global.Redis, err = redis.Init()
	if err != nil {
		global.Log.Printf("Init redis failed, err:%v\n", err)
		return
	}
	// 初始化 ES
	if global.Config.ES.Enable {
		global.ES, err = es.Init()
		if err != nil {
			global.Log.Errorf("es init err:%s\n", err.Error())
			return
		}
		// 创建文章索引
		a := model.ArticleModel{}
		err = a.CreateIndex()
		if err != nil {
			global.Log.Errorf("es index init err:%s\n", err.Error())
			return
		}
		// a.DeleteIndex()
	}

	// 解析命令行参数
	op := flag.FlagUserParse()

	// 如果没有传递用户名参数，启动服务器
	if op.Username == "" || op.Password == "" {
		startServer()
		return
	}

	// 如果有用户名参数，则尝试创建用户
	flag.CreateUser(&op)
}

// 启动 Gin 服务器
func startServer() {

	wg := sync.WaitGroup{}
	r := router.InitRouter(global.Config.System.Env, &wg)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", global.Config.System.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info("Shutdown Server ...")

	// 等待所有后台任务完成
	wg.Wait()

	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 优雅关闭服务
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Fatal("Server Shutdown: ", err.Error())
	}
	global.Log.Info("Server exiting")
}
