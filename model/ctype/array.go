package ctype

import (
	"database/sql/driver"
	"strings"
)

type ArrayString []string

// 读数据库的时候 string 变为 ArrayString 类型
func (a *ArrayString) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}

	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		*a = []string{}
		return nil
	}

	s = strings.TrimSpace(s)
	if s == "" {
		*a = []string{}
		return nil
	}

	s = strings.NewReplacer(
		"\r\n", "\n",
		"\r", "\n",
		"\n", ",",
		"，", ",",
	).Replace(s)

	parts := strings.Split(s, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}
	*a = res
	return nil
}

// Value 实现 driver.Valuer 接口，
// 向数据库中写数据的时候,ArrayString 转换为字符串,并且用","分开
func (a ArrayString) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	res := make([]string, 0, len(a))
	for _, v := range a {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		res = append(res, v)
	}
	return strings.Join(res, ","), nil
}

func (ArrayString) GormDataType() string {
	return "text"
}
