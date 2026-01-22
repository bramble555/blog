package jwt

import (
	"errors"
	"time"

	"github.com/bramble555/blog/global"
	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	SN                 int64  `json:"sn,string"` // Snowflake ID
	Role               int64  `json:"role"`
	Username           string `json:"user_name"`
	jwt.StandardClaims        // 嵌入 jwt.StandardClaims
}

// 不能在调用函数之前直接初始化，否则会空指针异常
// 因为函数执行顺序是 包级变量  函数

// GenToken 生成JWT
func GenToken(sn int64, role int64, username string) (string, error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		SN:       sn, // 自定义字段
		Role:     role,
		Username: username,
		StandardClaims: jwt.StandardClaims{ // 明确指定字段名
			ExpiresAt: time.Now().Add(time.Duration(global.Config.Jwt.Expries * int64(time.Hour))).Unix(),
			Issuer:    global.Config.Jwt.Issuer, // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(global.Config.Jwt.Secret))
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	secret := global.Config.Jwt.Secret
	if secret == "" {
		global.Log.Error("JWT Secret is empty in config!")
	}
	// 解析token
	var mc = new(MyClaims)
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		global.Log.Errorf("jwt.ParseWithClaims failed: %v", err)
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("token is invalid")
}
