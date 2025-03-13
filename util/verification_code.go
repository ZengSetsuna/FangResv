package util

import (
	"crypto/rand"
	"math/big"
)

// GenerateVerificationCode 生成 6 位随机数字验证码
func GenerateVerificationCode() string {
	const digits = "0123456789"
	length := 6
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			panic(err) // 生成失败，直接 panic（生产环境应返回错误）
		}
		code[i] = digits[n.Int64()]
	}
	return string(code)
}
