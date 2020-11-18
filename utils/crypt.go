package utils

import "golang.org/x/crypto/bcrypt"

// EncryptPassword 对密码进行加密
func EncryptPassword(password string) (string,error) {
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return "",err
	}
	return string(hashedPassword),nil
}

// CompareHashAndPassword 比较密码和加密后的哈希值
func CompareHashAndPassword(password,hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword),[]byte(password))
}
	