package config

type Email struct {
	Host             string `json:"host" mapstructure:"host"`
	Port             int    `json:"port" mapstructure:"port"`
	User             string `json:"user" mapstructure:"user"` //发件人邮箱
	Password         string `json:"password" mapstructure:"password"`
	DefaultFromEmail string `json:"default_from_email" mapstructure:"default_from_email"` //默认发件人名称
	UseSSL           bool   `json:"use_ssl" mapstructure:"use_ssl"`                       //是否使用ssl
	UserTLs          bool   `json:"user_tls" mapstructure:"user_tls"`
}
