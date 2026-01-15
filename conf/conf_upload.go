package config

type Upload struct {
	Size   int    `mapstructure:"size" json:"size"`       // 图片上传的大小
	Path   string `mapstructure:"path" json:"path"`       // Banner 图片上传目录
	AdPath string `mapstructure:"ad_path" json:"ad_path"` // Advert 图片上传目录
}
