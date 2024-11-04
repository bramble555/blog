package pkg

import (
	"errors"
	"time"

	"github.com/bramble555/blog/global"
	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	ID                 uint   `json:"id"`
	Role               uint   `json:"role"`
	Username           string `json:"user_name"`
	jwt.StandardClaims        // 嵌入 jwt.StandardClaims
}

// 不能在调用函数之前直接初始化，否则会空指针异常
// 因为函数执行顺序是 包级变量  函数
var tokeExpireDuration int64
var mySecret []byte

// GenToken 生成JWT
func GenToken(id uint, role uint, username string) (string, error) {
	tokeExpireDuration = global.Config.Jwt.Expries
	mySecret = []byte(global.Config.Jwt.Secret)
	// 创建一个我们自己的声明
	claims := MyClaims{
		ID:       id, // 自定义字段
		Role:     role,
		Username: username,
		StandardClaims: jwt.StandardClaims{ // 明确指定字段名
			ExpiresAt: time.Now().Add(time.Duration(tokeExpireDuration * int64(time.Hour))).Unix(),
			Issuer:    global.Config.Jwt.Issuer, // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		// 这里是验证签名的密钥，根据你的 JWT 签名算法来设置
		// 例如，如果你的 JWT 是使用 HS256 签名的，这里应该返回签名的密钥
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
