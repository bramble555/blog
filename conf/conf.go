package config

type Config struct {
	System        System    `mapstructure:"system"`
	Logger        Logger    `mapstructure:"logger"`
	Mysql         Mysql     `mapstructure:"mysql"`
	Redis         Redis     `mapstructure:"redis"`
	Email         Email     `mapstructure:"email"`
	Jwt           Jwt       `mapstructure:"jwt"`
	Upload        Upload    `mapstructure:"upload"`
	Elasticsearch ESConfig  `mapstructure:"elasticsearch"`
	RateLimit     RateLimit `mapstructure:"rate_limit"`
}
