package config

type Jwt struct {
	Secret  string `json:"secret" mapstructure:"secret"`   //密钥
	Expries int    `json:"expries" mapstructure:"expries"` //过期时间
	Issuer  string `json:"issuer" mapstructure:"issuer"`   //颁发人
}
