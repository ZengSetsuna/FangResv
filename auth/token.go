package auth

import (
	"errors"
	"time"

	"github.com/o1egl/paseto"
)

const secretKey = "dF453fEsEV3bjfnd29cFoLpq8432fn9O" // to be put into config file

func CreateToken(userID int32) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)

	token := paseto.NewV2()
	jsonToken := map[string]interface{}{
		"user_id": userID,
		"exp":     exp.Unix(),
	}

	return token.Encrypt([]byte(secretKey), jsonToken, nil)
}

func ValidateToken(token string) (map[string]interface{}, error) {
	var jsonToken map[string]interface{}
	err := paseto.NewV2().Decrypt(token, []byte(secretKey), &jsonToken, nil)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	return jsonToken, nil
}
