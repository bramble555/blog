package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bramble555/blog/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	dsn := global.Config.Mysql.Dsn()
	//打印mysql日志
	var mysqlLogger logger.Interface
	if global.Config.System.Env == "debug" {
		//开发环境显示所有的sql
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error) //只打印错误的sql
	}
	global.MysqlLog = mysqlLogger

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		global.Log.Fatalf(fmt.Sprintf("[%s]mysql连接失败", dsn))
		return nil, err
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	sqlDB.SetMaxIdleConns(10)               //最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              //最多可容纳
	sqlDB.SetConnMaxLifetime(time.Hour * 4) //连接最大复用时间，不能超过mysql的wait_timeout
	return db, err
}
