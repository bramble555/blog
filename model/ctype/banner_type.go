package ctype

import "encoding/json"

type BannerType int

const (
	Local BannerType = 1
	QiNiu BannerType = 2
)

// 解析成 json 格式的时候还是 string类型
func (s BannerType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s BannerType) String() string {
	var str string
	switch s {
	case Local:
		str = "本地"
	case QiNiu:
		str = "七牛云"
	default:
		str = "其他"
	}
	return str
}
