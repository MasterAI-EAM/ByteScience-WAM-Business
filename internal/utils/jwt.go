package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// GetToken 生成token
func GetToken(secretKey string, seconds int64, userId string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Unix() + seconds
	claims["iat"] = time.Now().Unix()
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
