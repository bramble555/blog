package pkg

import "golang.org/x/crypto/bcrypt"

// HashPassword 生成的哈希密码是59~72位，当然也不能接收超越72位的密码，所以需要数据库varchar(72)
func HashPassword(password string) (string, error) {
	// 使用 bcrypt 库的 GenerateFromPassword 函数进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
func ComparePasswords(hashedPassword, inputPassword string) error {
	// 使用 bcrypt 库的 CompareHashAndPassword 函数比较密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		
	}
	return err
}
