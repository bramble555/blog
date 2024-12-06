package config

type Config struct {
	System   System   `mapstructure:"system"`
	Logger   Logger   `mapstructure:"logger"`
	Mysql    Mysql    `mapstructure:"mysql"`
	Redis    Redis    `mapstructure:"redis"`
	SiteInfo SiteInfo `mapstructure:"site_info"`
	QQ       QQ       `mapstructure:"qq"`
	Qiniu    Qiniu    `mapstructure:"qi_niu"`
	Email    Email    `mapstructure:"email"`
	Jwt      Jwt      `mapstructure:"jwt"`
	Upload   Upload   `mapstructure:"upload"`
	ES       ES       `mapstructure:"es"`
}
