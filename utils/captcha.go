package utils

import (
	"crypto/rand"
)

var table = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

// GenerateCode 生成验证码
func GenerateCode(max int) string {
	code := make([]byte, max)
	rand.Read(code)

	for i, v := range code {
		code[i] = table[int(v)%len(table)]
	}

	return string(code)
}
