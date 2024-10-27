package ctype

import (
	"strconv"
	"strings"

	"github.com/bramble555/blog/global"
)

type NumString string

func (ns *NumString) Scan(s string) error {
	_, err := strconv.Atoi(s)
	if strings.ToLower(s) == "null " || err == nil {
		return nil
	}
	global.Log.Errorf("numString error:%s\n", err.Error())
	return err
}
