package config

type Logger struct {
	Level    string `mapstructure:"level"`
	Prefix   string `mapstructure:"prefix"`
	Director string `mapstructure:"director"`
	ShowLine bool   `mapstructure:"show_line"` //显示行号
	FileName string `mapstructure:"file_name"`
}
