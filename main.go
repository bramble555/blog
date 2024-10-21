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
	// 注册路由
	r := router.InitRounter(global.Config.System.Env)
	// 优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", global.Config.System.Port),
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	global.Log.Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Fatal("Server Shutdown: ", err.Error())
	}
	global.Log.Info("Server exiting")
}
