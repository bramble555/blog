package config

import "strconv"

type Mysql struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Config   string `mapstructure:"config"` //高级配置
	DB       string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	LogLevel string `mapstructure:"log_level"` //日志等级，debu\dev\release
}

func (m *Mysql) Dsn() string {
	return m.User + ":" + m.Password +
		"@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) +
		")/" + m.DB + "?" + m.Config
}
