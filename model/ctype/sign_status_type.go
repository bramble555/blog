package ctype

import "encoding/json"

type SignStatus int

const (
	SignQQ SignStatus = 1 // QQ
)

func (s SignStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s SignStatus) String() string {
	var str string
	switch s {
	case SignQQ:
		str = "QQ"
	default:
		str = "其他"
	}
	return str
}
