package config

type ES struct {
	Enable bool   `json:"enable" mapstructure:"enable"`
	Host   string `json:"host" mapstructure:"host"`
	Port   int    `json:"port" mapstructure:"port"`
	Sniff  bool   `json:"sniff" mapstructure:"sniff"`
}
