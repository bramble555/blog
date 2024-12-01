package ctype

import "encoding/json"

type SignStatus int

const (
	SignQQ     SignStatus = 1 // QQ
	SignGithub SignStatus = 2 // Github
	SignEmail  SignStatus = 3 // 邮箱
)

func (s SignStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s SignStatus) String() string {
	var str string
	switch s {
	case SignQQ:
		str = "QQ"
	case SignGithub:
		str = "Github"
	case SignEmail:
		str = "邮箱"
	default:
		str = "其他"
	}
	return str
}
