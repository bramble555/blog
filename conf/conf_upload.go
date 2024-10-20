package config

type Upload struct {
	Size int    `mapstructure:"size" json:"size"` //图片上传的大小
	Path string `mapstructure:"path" json:"path"` //图片上传的目录
}
