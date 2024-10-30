package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bramble555/blog/dao/mysql"
	"github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/flag"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logger"
	"github.com/bramble555/blog/router"
	"github.com/bramble555/blog/setting"
)

func main() {
	var err error
	// 初始化config
	if err = setting.Init(); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
		return
	}
	// 初始化logger
	if global.Log, err = logger.Init(); err != nil {
		fmt.Printf("Init logger failed, err:%v\n", err)
		return
	}
	// 初始化db
	if global.DB, err = mysql.Init(); err != nil {
		global.Log.Errorf("Init mysql failed, err:%v\n", err)
		return
	}
	// 初始化redis
	global.Redis, err = redis.Init()
	if err != nil {
		global.Log.Printf("Init redis failed, err:%v\n", err)
		return
	}

	// 解析命令行参数
	op := flag.FlagUserParse()

	// 如果没有传递用户名参数，启动服务器
	if op.Username == "" {
		startServer()
		return
	}

	// 如果有用户名参数，则尝试创建用户
	flag.CreateUser()
}

// 启动 Gin 服务器
func startServer() {
	r := router.InitRounter(global.Config.System.Env)
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

	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Fatal("Server Shutdown: ", err.Error())
	}
	global.Log.Info("Server exiting")
}
