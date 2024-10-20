package config

type Qiniu struct {
	Enable    bool    `json:"enable" mapstructure:"enable"` //是否启用
	AccessKey string  `json:"access_key" mapstructure:"access_key"`
	SecretKey string  `json:"secret_key" mapstructure:"secret_key"`
	Bucket    string  `json:"bucket" mapstructure:"bucket"` //存储桶的名字
	CDN       string  `json:"cdn" mapstructure:"cdn"`       //访问图片的地址前缀
	Prefix    string  `json:"prefix" mapstructure:"prefix"` //前缀
	Zone      string  `json:"zone" mapstructure:"zone"`     //存储的地区
	Size      float64 `json:"size" mapstructure:"size"`     //图片大小限制，单位为MB
}
