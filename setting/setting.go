package setting

import (
	"fmt"

	"github.com/bramble555/blog/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("setting")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&global.Config)
	// 监听配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改了")
	})
	return
}
