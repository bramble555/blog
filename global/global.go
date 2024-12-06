package global

import (
	config "github.com/bramble555/blog/conf"
	"github.com/olivere/elastic/v7"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Config   *config.Config // 配置文件
	DB       *gorm.DB       // 数据库文件
	Log      *logrus.Logger // 日志文件
	MysqlLog logger.Interface
	Redis    *redis.Client
	ES       *elastic.Client
)
