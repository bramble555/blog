package config

type Email struct {
	Host             string `json:"host" mapstructure:"host"`
	Port             int    `json:"port" mapstructure:"port"`
	User             string `json:"user" mapstructure:"user"`
	Password         string `json:"password" mapstructure:"password"`
	DefaultFromEmail string `json:"default_from_email" mapstructure:"default_from_email"`
}
  