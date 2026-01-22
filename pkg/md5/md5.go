package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 实现加密功能
func MD5(src []byte) string {
	sum := md5.Sum(src)
	return hex.EncodeToString(sum[:])
}
