package config

import "fmt"

type QQ struct {
	AppID    string `json:"app_id" mapstructure:"app_id"`
	Key      string `json:"key" mapstructure:"key"`
	Redirect string `json:"redirect" mapstructure:"redirect"` //登陆后的回调地址
}

func (q QQ) GetPath() string {
	if q.Key == "" || q.AppID == "" || q.Redirect == "" {
		return ""
	}
	return fmt.Sprintf("%s+%s", q.AppID, q.Redirect)
}
