package auth

import (
	"errors"
	"time"

	"github.com/o1egl/paseto"
)

// 定义 TokenPayload 结构体
type TokenPayload struct {
	UserID int32 `json:"user_id"`
	Exp    int64 `json:"exp"`
}

// 定义常量，表示不同的状态
const (
	TokenOK      = iota // 0: token 有效
	TokenExpired        // 1: token 已过期
	TokenInvalid        // 2: token 无效
)

// 定义密钥
const secretKey = "dF453fEsEV3bjfnd29cFoLpq8432fn9O" // to be put into config file

// CreateToken 生成一个有效期为 24 小时的 token
func CreateToken(userID int32) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)

	token := paseto.NewV2()

	// 使用结构体代替 map
	payload := TokenPayload{
		UserID: userID,
		Exp:    exp.Unix(),
	}

	// 使用结构体来创建 token
	return token.Encrypt([]byte(secretKey), payload, nil)
}

// ValidateToken 验证 token 并解析出结构体
// 返回 token 的状态：有效、过期或无效，以及对应的 payload
func ValidateToken(token string) (int, *TokenPayload, error) {
	var payload TokenPayload

	// 解密并填充结构体
	err := paseto.NewV2().Decrypt(token, []byte(secretKey), &payload, nil)
	if err != nil {
		return TokenInvalid, nil, errors.New("invalid token") // 无效 token
	}

	// 检查 token 是否过期
	if time.Now().Unix() > payload.Exp {
		return TokenExpired, nil, errors.New("token expired") // 过期 token
	}

	// 如果 token 有效，返回有效状态
	return TokenOK, &payload, nil
}
